package entity

import "time"

type Message struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Content    string    `json:"content"`
	Type       string    `json:"type"`
	Timestamp  int64     `json:"timestamp"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

const (
	MessageTypePrivate = "private"
	MessageTypeGroup   = "group"
)

type OfflineMessage struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	MessageID string    `json:"message_id"`
	CreatedAt time.Time `json:"created_at"`
}