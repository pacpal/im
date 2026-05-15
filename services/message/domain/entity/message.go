package entity

import "time"

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

type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeVideo    MessageType = "video"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeFile     MessageType = "file"
	MessageTypeLocation MessageType = "location"
)

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

func (m *Message) MarkAsRead() {
	m.IsRead = true
	m.ReadAt = time.Now().UnixMilli()
}

func (m *Message) Revoke() {
	m.IsRevoked = true
}

func (m *Message) IsGroupMessage() bool {
	return m.MsgType == MessageTypeText
}

type OfflineMessage struct {
	UserID   string
	Messages []*Message
}
