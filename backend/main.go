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

type SentimentRequest struct {
	Sentence string `json:"sentence"`
}

type SentimentResponse struct {
	PredictedLabel string  `json:"predicted_label"`
	Confidence     float64 `json:"confidence"`
}

func main() {
	db.InitDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/predict", handleSentAnalysis)
	mux.HandleFunc("/login", handlers.HandleLogin)
	mux.HandleFunc("/signup", handlers.HandleSignup)
	mux.HandleFunc("/dashboard/saveEmotion", handlers.SaveEmotionHandler)
	mux.HandleFunc("/dashboard/getEmotions", handlers.GetEmotionsHandler)
	mux.HandleFunc("/dashboard", middleware.AuthMiddleware(handlers.HandleDashboard))

	// Wrap the mux with CORS middleware
	corsHandler := middleware.CORSMiddleware(mux)

	fmt.Println("Go server is running on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

func handleSentAnalysis(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	fmt.Println("Request Body:", string(body))

	// Forward the Request to the Python backend
	pythonBackendURL := "http://localhost:8000/predict"
	resp, err := http.Post(pythonBackendURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error connecting to Python backend: %v\n", err)
		http.Error(w, "Failed to connect to Python backend", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the Python backend response
	pythonResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read Python backend response", http.StatusInternalServerError)
		return
	}
	fmt.Println("Python Backend Response:", string(pythonResponse))

	// Forward the Python backend response to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(pythonResponse)
}
