package service

import (
	"IM/pkg/id"
	"IM/services/message/domain/entity"
	"IM/services/message/domain/event"
	"IM/services/message/domain/repository"
	"context"
	"errors"
	"time"
)

var (
	ErrMessageNotFound = errors.New("message not found")
	ErrNotSender       = errors.New("not message sender")
	ErrAlreadyRevoked  = errors.New("message already revoked")
)

type MessageService struct {
	messageRepo     repository.MessageRepository
	messageCache    repository.MessageCache
	messageProducer MessageProducer
	idGenerator     *id.SnowflakeGenerator
	eventPublisher  *event.EventPublisher
}

type MessageProducer interface {
	PublishMessage(ctx context.Context, msg *entity.Message) error
}

func NewMessageService(
	messageRepo repository.MessageRepository,
	messageCache repository.MessageCache,
	messageProducer MessageProducer,
	idGenerator *id.SnowflakeGenerator,
	eventPublisher *event.EventPublisher,
) *MessageService {
	return &MessageService{
		messageRepo:     messageRepo,
		messageCache:    messageCache,
		messageProducer: messageProducer,
		idGenerator:     idGenerator,
		eventPublisher:  eventPublisher,
	}
}

func (s *MessageService) SendMessage(ctx context.Context, senderID, receiverID, content, msgType string) (*entity.Message, error) {
	messageID := s.idGenerator.Generate()
	message := entity.NewMessage(messageID, senderID, receiverID, content, entity.MessageType(msgType))

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	if err := s.messageCache.IncrUnreadCount(ctx, receiverID); err != nil {
	}

	if err := s.messageProducer.PublishMessage(ctx, message); err != nil {
	}

	s.eventPublisher.Publish(&event.MessageSentEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "message.sent",
			OccurredAt:  time.Now(),
			AggregateID: messageID,
		},
		MessageID:  messageID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		MsgType:    msgType,
	})

	return message, nil
}

func (s *MessageService) GetMessage(ctx context.Context, messageID string) (*entity.Message, error) {
	return s.messageRepo.GetByID(ctx, messageID)
}

func (s *MessageService) GetOfflineMessages(ctx context.Context, userID string, limit, offset int) ([]*entity.Message, error) {
	return s.messageRepo.GetByReceiverID(ctx, userID, limit, offset)
}

func (s *MessageService) GetHistoryMessages(ctx context.Context, userID, targetID string, beforeTime int64, limit int) ([]*entity.Message, error) {
	return s.messageRepo.GetHistory(ctx, userID, targetID, beforeTime, limit)
}

func (s *MessageService) MarkAsRead(ctx context.Context, messageID, userID string) error {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return ErrMessageNotFound
	}

	if message.ReceiverID != userID {
		return errors.New("not message receiver")
	}

	if message.IsRead {
		return nil
	}

	if err := s.messageRepo.MarkAsRead(ctx, messageID); err != nil {
		return err
	}

	if err := s.messageCache.DecrUnreadCount(ctx, userID); err != nil {
	}

	s.eventPublisher.Publish(&event.MessageReadEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "message.read",
			OccurredAt:  time.Now(),
			AggregateID: messageID,
		},
		MessageID: messageID,
		UserID:    userID,
	})

	return nil
}

func (s *MessageService) MarkAllAsRead(ctx context.Context, userID, targetID string) error {
	return s.messageRepo.MarkAllAsRead(ctx, userID, targetID)
}

func (s *MessageService) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	count, err := s.messageCache.GetUnreadCount(ctx, userID)
	if err != nil {
		return s.messageRepo.GetUnreadCount(ctx, userID)
	}
	return count, nil
}

func (s *MessageService) GetOnlineStatus(ctx context.Context, userIDs []string) (map[string]bool, error) {
	onlineUsers, err := s.messageCache.GetOnlineUsers(ctx)
	if err != nil {
		return make(map[string]bool), nil
	}

	result := make(map[string]bool)
	for _, userID := range userIDs {
		result[userID] = onlineUsers[userID]
	}
	return result, nil
}

func (s *MessageService) RevokeMessage(ctx context.Context, messageID, userID string) error {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return ErrMessageNotFound
	}

	if message.SenderID != userID {
		return ErrNotSender
	}

	if message.IsRevoked {
		return ErrAlreadyRevoked
	}

	if err := s.messageRepo.Revoke(ctx, messageID); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.MessageRevokedEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "message.revoked",
			OccurredAt:  time.Now(),
			AggregateID: messageID,
		},
		MessageID: messageID,
		UserID:    userID,
	})

	return nil
}

func (s *MessageService) SetUserOnline(ctx context.Context, userID string) error {
	if err := s.messageCache.SetOnlineUser(ctx, userID); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.UserOnlineEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "user.online",
			OccurredAt:  time.Now(),
			AggregateID: userID,
		},
		UserID: userID,
	})

	return nil
}

func (s *MessageService) SetUserOffline(ctx context.Context, userID string) error {
	if err := s.messageCache.RemoveOnlineUser(ctx, userID); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.UserOfflineEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "user.offline",
			OccurredAt:  time.Now(),
			AggregateID: userID,
		},
		UserID: userID,
	})

	return nil
}
