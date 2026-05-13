package service

import (
	"IM/services/message-service/domain/entity"
	"IM/services/message-service/infrastructure/repository"
	ws "IM/services/message-service/infrastructure/websocket"
	"context"
	"errors"
	"time"
)

var (
	ErrMessageNotFound     = errors.New("message not found")
	ErrNotFriends          = errors.New("not friends")
	ErrNotGroupMember      = errors.New("not a group member")
	ErrInvalidMessage      = errors.New("invalid message")
)

type MessageApplicationService struct {
	messageRepo    repository.MessageRepository
	friendshipRepo repository.FriendshipRepository
	groupMemberRepo repository.GroupMemberRepository
	hub            *ws.Hub
}

func NewMessageApplicationService(
	messageRepo repository.MessageRepository,
	friendshipRepo repository.FriendshipRepository,
	groupMemberRepo repository.GroupMemberRepository,
) *MessageApplicationService {
	return &MessageApplicationService{
		messageRepo:    messageRepo,
		friendshipRepo: friendshipRepo,
		groupMemberRepo: groupMemberRepo,
		hub:            ws.NewHub(),
	}
}

func (s *MessageApplicationService) GetHub() *ws.Hub {
	return s.hub
}

func (s *MessageApplicationService) StartHub() {
	go s.hub.Run()
}

func (s *MessageApplicationService) SendMessage(ctx context.Context, senderID, receiverID, content, msgType string) (*entity.Message, error) {
	if content == "" {
		return nil, ErrInvalidMessage
	}

	msg := &entity.Message{
		ID:         generateMessageID(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Type:       msgType,
		Timestamp:  time.Now().UnixMilli(),
		IsRead:     false,
		CreatedAt:  time.Now(),
	}

	switch msgType {
	case entity.MessageTypePrivate:
		return msg, s.sendPrivateMessage(ctx, msg)
	case entity.MessageTypeGroup:
		return msg, s.sendGroupMessage(ctx, msg)
	default:
		return nil, ErrInvalidMessage
	}
}

func (s *MessageApplicationService) sendPrivateMessage(ctx context.Context, msg *entity.Message) error {
	isFriend, err := s.friendshipRepo.Exists(ctx, msg.SenderID, msg.ReceiverID)
	if err != nil {
		return err
	}
	if !isFriend {
		return ErrNotFriends
	}

	if err := s.messageRepo.Create(ctx, msg); err != nil {
		return err
	}

	s.hub.BroadcastMessage(msg)
	return nil
}

func (s *MessageApplicationService) sendGroupMessage(ctx context.Context, msg *entity.Message) error {
	isMember, err := s.groupMemberRepo.IsMember(ctx, msg.ReceiverID, msg.SenderID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrNotGroupMember
	}

	if err := s.messageRepo.Create(ctx, msg); err != nil {
		return err
	}

	memberIDs, err := s.groupMemberRepo.GetMemberIDs(ctx, msg.ReceiverID)
	if err != nil {
		return err
	}

	for _, memberID := range memberIDs {
		if memberID != msg.SenderID {
			msgCopy := *msg
			msgCopy.ReceiverID = memberID
			s.hub.BroadcastMessage(&msgCopy)
		}
	}

	return nil
}

func (s *MessageApplicationService) GetMessage(ctx context.Context, messageID string) (*entity.Message, error) {
	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, ErrMessageNotFound
	}
	return msg, nil
}

func (s *MessageApplicationService) GetOfflineMessages(ctx context.Context, userID string, limit, offset int) ([]*entity.Message, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.messageRepo.GetOfflineMessages(ctx, userID, limit, offset)
}

func (s *MessageApplicationService) MarkAsRead(ctx context.Context, messageID string) error {
	return s.messageRepo.MarkAsRead(ctx, messageID)
}

func (s *MessageApplicationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.messageRepo.MarkAllAsRead(ctx, userID)
}

func (s *MessageApplicationService) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return s.messageRepo.GetUnreadCount(ctx, userID)
}

func (s *MessageApplicationService) GetChatHistory(ctx context.Context, userID1, userID2 string, limit, offset int) ([]*entity.Message, error) {
	if limit <= 0 {
		limit = 50
	}

	var allMessages []*entity.Message

	sent, err := s.messageRepo.GetBySender(ctx, userID1, limit, offset)
	if err != nil {
		return nil, err
	}

	received, err := s.messageRepo.GetByReceiver(ctx, userID1, limit, offset)
	if err != nil {
		return nil, err
	}

	allMessages = append(allMessages, sent...)
	allMessages = append(allMessages, received...)

	return allMessages, nil
}

func generateMessageID() string {
	return "msg_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}