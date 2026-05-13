package entity

// Message 消息领域实体
type Message struct {
	ID         string
	SenderID   string
	ReceiverID string
	Content    string
	Type       string // "private" 或 "group"
	Timestamp  int64
	IsRead     bool
	CreatedAt  int64
}
