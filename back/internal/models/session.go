package models

type Session struct {
	SessionID string `json:"session_id" redis:"session_id"`
	UserID    int    `json:"user_id" redis:"user_id"`
}
