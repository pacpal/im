// Package entity 定义消息领域实体。
package entity

import "time"

// Message 表示一条消息实体。
type Message struct {
	ID         string
	SenderID   string
	ReceiverID string
	Content    string
	MsgType    MessageType
	Timestamp  int64
	IsRead     bool
	IsRevoked  bool
	ReadAt     int64
	CreatedAt  time.Time
}

// MessageType 枚举消息类型。
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeVideo    MessageType = "video"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeFile     MessageType = "file"
	MessageTypeLocation MessageType = "location"
)

// NewMessage 创建消息实体并填充时间戳。
func NewMessage(id, senderID, receiverID, content string, msgType MessageType) *Message {
	now := time.Now()
	return &Message{
		ID:         id,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		MsgType:    msgType,
		Timestamp:  now.UnixMilli(),
		IsRead:     false,
		IsRevoked:  false,
		CreatedAt:  now,
	}
}

// MarkAsRead 标记消息为已读并记录时间。
func (m *Message) MarkAsRead() {
	m.IsRead = true
	m.ReadAt = time.Now().UnixMilli()
}

// Revoke 将消息标记为已撤回。
func (m *Message) Revoke() {
	m.IsRevoked = true
}

// IsGroupMessage 判断消息是否为群消息（示例实现）。
func (m *Message) IsGroupMessage() bool {
	return m.MsgType == MessageTypeText
}

// OfflineMessage 用于批量返回用户的离线消息。
type OfflineMessage struct {
	UserID   string
	Messages []*Message
}
