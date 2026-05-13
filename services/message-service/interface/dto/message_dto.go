package dto

import "IM/services/message-service/domain/entity"

type SendMessageRequest struct {
	ReceiverID string `json:"receiver_id"`
	Content   string `json:"content"`
	Type      string `json:"type"`
}

type MessageResponse struct {
	ID         string `json:"id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	Timestamp  int64  `json:"timestamp"`
	IsRead     bool   `json:"is_read"`
}

type GetMessagesRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func ToMessageResponse(msg *entity.Message) *MessageResponse {
	if msg == nil {
		return nil
	}
	return &MessageResponse{
		ID:         msg.ID,
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		Type:       msg.Type,
		Timestamp:  msg.Timestamp,
		IsRead:     msg.IsRead,
	}
}

func ToMessageResponseList(msgs []*entity.Message) []*MessageResponse {
	result := make([]*MessageResponse, len(msgs))
	for i, m := range msgs {
		result[i] = ToMessageResponse(m)
	}
	return result
}