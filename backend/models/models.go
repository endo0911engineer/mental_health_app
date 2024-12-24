package models

import "time"

type Emotion struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Emotion   string    `json:"emotion"`
	Score     int       `json:"score"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"createdAt"`
}

// ユーザー情報
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmotionAnalysisResult struct {
	Score int `json:"score"`
}

type AnalysisRequest struct {
	Text string `json:"text"`
}

type AnalysisResponse struct {
	Keywords  []string            `json:"keywords"`
	Relations map[string][]string `json:"relations"`
}
