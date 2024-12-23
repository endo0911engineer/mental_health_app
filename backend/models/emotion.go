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
