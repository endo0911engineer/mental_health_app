package middleware

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

// シークレットキー
var mySecretKey = 1234
var JwtSecret = []byte(strconv.Itoa(mySecretKey))

// JWTクレーム構造体
type Claims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// AuthMiddleware verifies the JWT token
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
