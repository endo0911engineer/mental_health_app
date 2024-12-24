package db

import (
	"backend/models"
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./backend.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	// データベース初期化
	err = createTable()
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
}

func createTable() error {
	// usersテーブル作成
	userQuery := `
		CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
		);
		`

	_, err := db.Exec(userQuery)
	if err != nil {
		return err
	}

	// emotionテーブル作成
	emotionQuery := `
	CREATE TABLE IF NOT EXISTS emotions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	emotion_text TEXT NOT NULL,
	score INTEGER NOT NULL,
	date DATE NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`

	_, err = db.Exec(emotionQuery)
	return err
}

// 新しいユーザーを作成
func CreateUser(user models.User) error {
	query := `INSERT INTO users (email, password) VALUES (?, ?)`
	_, err := db.Exec(query, user.Email, user.Password)
	return err
}

// メールアドレスでユーザーを取得
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password FROM users WHERE email = ?`
	row := db.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// ユーザーが見つからない場合
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// 感情をデータベースに保存
func SaveEmotion(emotion models.Emotion) error {
	query := `
	INSERT INTO emotions (user_id, emotion_text, score, date)
	VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, emotion.UserID, emotion.Emotion, emotion.Score, emotion.Date)
	return err
}

// 更新された感情をデータベースに保存
func UpdateEmotion(emotion models.Emotion) error {
	query := `
	UPDATE emotions
	SET emotion = ?, score = ? 
	WHERE user_id = ? AND date = ?
	`
	_, err := db.Exec(query, emotion.Emotion, emotion.Score, emotion.UserID, emotion.Date)
	return err
}

// 感情をデータベースから削除
func DeleteEmotion(userID int, date time.Time) error {
	query := `
	DELETE FROM emotions
	WHERE user_id = ? AND date = ?
	`
	_, err := db.Exec(query, userID, date)
	return err
}

// ユーザーIDと日付で感情を取得
func GetEmotionsByUserID(userID int) ([]models.Emotion, error) {
	var emotions []models.Emotion
	query := `SELECT id, user_id, emotion_text, score, date, created_at FROM emotions WHERE user_id = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var emotion models.Emotion
		err := rows.Scan(&emotion.ID, &emotion.UserID, &emotion.Emotion, &emotion.Score, &emotion.Date, &emotion.CreatedAt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		emotions = append(emotions, emotion)
	}

	return emotions, nil
}

// ユーザーIDに基づいて過去の感情データを取得
func GetPastEmotions(userID int) ([]models.Emotion, error) {
	var emotions []models.Emotion
	query := `SELECT emotion_text FROM emotions WHERE user_id = ? ORDER BY date DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var emotion models.Emotion
		err := rows.Scan(&emotion.ID, &emotion.UserID, &emotion.Emotion, &emotion.Score, &emotion.Date, &emotion.CreatedAt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		emotions = append(emotions, emotion)
	}

	return emotions, nil
}
