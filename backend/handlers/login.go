package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"backend/db"
	"backend/middleware"
	"backend/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// パスワード照合用関数
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HandleLogin handlers user Login and JWT generation
func HandleLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// データベースからユーザー情報を取得
	dbUser, err := db.GetUserByEmail(user.Email)
	if err != nil || dbUser == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// パスワードの照合
	if !checkPasswordHash(user.Password, dbUser.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// トークン生成
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &middleware.Claims{
		ID:    dbUser.ID,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtSecret)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"token":  tokenString,
		"userId": claims.ID,
	}

	log.Println("Login successful for user:", dbUser.Email)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
