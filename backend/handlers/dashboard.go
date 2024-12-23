package handlers

import (
	"backend/db"
	"backend/middleware"
	"backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*middleware.Claims)
	if !ok || claims == nil {
		http.Error(w, "Umauthorized access", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome to the dashboard",
		"user_id": strconv.Itoa(claims.ID),
	})
}

// HandleDashboard handles requests to the dashboard
func SaveEmotionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// クエリパラメータから user_id を取得
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		log.Printf("Missing user_id in query parameters")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// リクエストボディを読み取ってログに出力
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// ボディの内容をログに出力
	fmt.Println("Request Body:", string(body))

	var emotion models.Emotion
	err = json.Unmarshal(body, &emotion)
	if err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("Invalid user_id: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	// クエリパラメータから取得したuser_idをemotionにセット
	emotion.UserID = userIDInt

	// Python APIで感情スコアを分析
	score, err := analyzeEmotionWithAI(emotion.Emotion)
	if err != nil {
		log.Printf("Failed to analyze emotion: %v", err)
		http.Error(w, "Failed to analyze emotion", http.StatusInternalServerError)
		return
	}

	// 分析したスコアをEmotionにセット
	emotion.Score = int(score.Score)

	// 感情をデーターベースに保存
	err = db.SaveEmotion(emotion)
	if err != nil {
		log.Printf("Failed to save emotion to database: %v", err)
		http.Error(w, "Failed to save emotion", http.StatusInternalServerError)
		return
	}

	log.Printf("Emotion saved successfully: %+v", emotion)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Emotion saved successfully",
	})
}

// UpdateEmotionHandler handles requests to update an emotion
func UpdateEmotionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// クエリパラメータから user_id を取得
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		log.Printf("Missing user_id in query parameters")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// リクエストボディを読み取ってログに出力
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// ボディの内容をログに出力
	fmt.Println("Request Body:", string(body))

	var emotion models.Emotion
	err = json.Unmarshal(body, &emotion)
	if err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("Invalid user_id: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	// クエリパラメータから取得したuser_idをemotionにセット
	emotion.UserID = userIDInt

	// Python APIで感情スコアを分析
	score, err := analyzeEmotionWithAI(emotion.Emotion)
	if err != nil {
		log.Printf("Failed to analyze emotion: %v", err)
		http.Error(w, "Failed to analyze emotion", http.StatusInternalServerError)
		return
	}

	// 分析したスコアをEmotionにセット
	emotion.Score = int(score.Score)

	// 感情をデーターベースに保存
	err = db.SaveEmotion(emotion)
	if err != nil {
		log.Printf("Failed to save emotion to database: %v", err)
		http.Error(w, "Failed to save emotion", http.StatusInternalServerError)
		return
	}

	log.Printf("Emotion saved successfully: %+v", emotion)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Emotion saved successfully",
	})
}

// getEmotionsHandler handles requests to get emotions by user ID and date
func GetEmotionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	// user_idをintに変換
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id parameter", http.StatusBadRequest)
		return
	}

	emotions, err := db.GetEmotionsByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve emotions", http.StatusInternalServerError)
		return
	}

	// エンコードして返す前にログを出力
	log.Printf("Emotions retrieved: %v", emotions)

	// 正常に取得できたらエンコードしてレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(emotions)
	if err != nil {
		log.Printf("Error encoding response: %v", err) // エンコードエラーのログ
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteEmotionHandler handles requests to delete an emotion by user ID and date
func DeleteEmotionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// リクエストからuser_idを取得
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	// user_idをintに変換
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id parameter", http.StatusBadRequest)
		return
	}

	// リクエストからdateを取得
	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, "Missing date parameter", http.StatusBadRequest)
		return
	}

	dateTim, err := time.Parse("2006-01-02", date)
	if err != nil {
		http.Error(w, "Invalid date parameter", http.StatusBadRequest)
		return
	}

	// 感情データを削除
	err = db.DeleteEmotion(userID, dateTim)
	if err != nil {
		http.Error(w, "Failed to delete emotion", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Emotion deleted successfully"))
}

// Python APIを呼び出して感情分析を行う
func analyzeEmotionWithAI(emotion string) (*EmotionAnalysisResult, error) {
	// PythonのAPIエンドポイント
	url := "http://localhost:8000/predict"

	// リクエストボディに分析するテキストを設定
	requestBody, err := json.Marshal(map[string]string{
		"emotion": emotion,
	})
	if err != nil {
		return nil, err
	}

	// PythonのAPIへPOSTリクエストを送る
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// レスポンスの処理
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to analyze emotion, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var result EmotionAnalysisResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type EmotionAnalysisResult struct {
	Score int `json:"score"`
}
