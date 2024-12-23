package main

import (
	"backend/db"
	"backend/handlers"
	"backend/middleware"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	db.InitDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/login", handlers.HandleLogin)
	mux.HandleFunc("/signup", handlers.HandleSignup)
	mux.HandleFunc("/dashboard/saveEmotion", handlers.SaveEmotionHandler)
	mux.HandleFunc("/dashboard/updateEmotion", handlers.UpdateEmotionHandler)
	mux.HandleFunc("/dashboard/getEmotions", handlers.GetEmotionsHandler)
	mux.HandleFunc("/dashboard/deleteEmotion", handlers.DeleteEmotionHandler)
	mux.HandleFunc("/dashboard", middleware.AuthMiddleware(handlers.HandleDashboard))

	// Wrap the mux with CORS middleware
	corsHandler := middleware.CORSMiddleware(mux)

	fmt.Println("Go server is running on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

type AnalysisRequest struct {
	Text string `json:"text"`
}

type AnalysisResponse struct {
	Keywords  []string            `json:"keywords"`
	Relations map[string][]string `json:"relations"`
}

func AnalyzeTextHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// フロントエンドから日記テキストを受け取る
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Python APIに転送する
	pythonAPI := "http://localhost:5000/analyze" // Pythonのエンドポイント
	resp, err := http.Post(pythonAPI, "application/json", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Failed to connect to Python service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Pythonからのレスポンスをフロントに返す
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read Python response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}
