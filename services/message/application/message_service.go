// Package service 提供 message 服务的业务逻辑实现。
package service

import (
	"IM/pkg/id"
	"IM/pkg/logger"
	pkgevent "IM/pkg/event"
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

// MessageService 封装消息发送、获取、标记已读等业务行为。
type MessageService struct {
	messageRepo     repository.MessageRepository
	messageCache    repository.MessageCache
	messageProducer MessageProducer
	idGenerator     *id.SnowflakeGenerator
	eventPublisher  pkgevent.Publisher
}

// MessageProducer 定义消息发布的接口（MQ 层）。
type MessageProducer interface {
	PublishMessage(ctx context.Context, msg *entity.Message) error
}

// NewMessageService 创建 MessageService 实例。
func NewMessageService(
	messageRepo repository.MessageRepository,
	messageCache repository.MessageCache,
	messageProducer MessageProducer,
	idGenerator *id.SnowflakeGenerator,
	eventPublisher pkgevent.Publisher,
) *MessageService {
	return &MessageService{
		messageRepo:     messageRepo,
		messageCache:    messageCache,
		messageProducer: messageProducer,
		idGenerator:     idGenerator,
		eventPublisher:  eventPublisher,
	}
}

func (s *MessageService) SendMessage(ctx context.Context, senderID, receiverID, content, msgType string) (res *entity.Message, err error) {
	done := logger.StartStep("MessageService.SendMessage", "sender", senderID, "receiver", receiverID)
	defer func() { done(err) }()

	messageID := s.idGenerator.Generate()
	res = entity.NewMessage(messageID, senderID, receiverID, content, entity.MessageType(msgType))

	if err = s.messageRepo.Create(ctx, res); err != nil {
		return nil, err
	}

	if e := s.messageCache.IncrUnreadCount(ctx, receiverID); e != nil {
		logger.Warnw("SendMessage: incr unread failed", "component", "message_service", "err", e, "receiver", receiverID)
	}

	if e := s.messageProducer.PublishMessage(ctx, res); e != nil {
		logger.Warnw("SendMessage: publish message failed", "component", "message_service", "err", e, "msg_id", messageID)
	}

	s.eventPublisher.Publish(&event.MessageSentEvent{
		BaseEvent: pkgevent.BaseEvent{
			EventType:   "message.sent",
			OccurredAt:  time.Now(),
			AggregateID: messageID,
		},
		MessageID:  messageID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		MsgType:    msgType,
	})

	logger.Infow("SendMessage: sent", "component", "message_service", "msg_id", messageID)
	return res, nil
}

func (s *MessageService) GetMessage(ctx context.Context, messageID string) (res *entity.Message, err error) {
	done := logger.StartStep("MessageService.GetMessage", "msg_id", messageID)
	defer func() { done(err) }()

	res, err = s.messageRepo.GetByID(ctx, messageID)
	return
}

func (s *MessageService) GetOfflineMessages(ctx context.Context, userID string, limit, offset int) (res []*entity.Message, err error) {
	done := logger.StartStep("MessageService.GetOfflineMessages", "user_id", userID, "limit", limit, "offset", offset)
	defer func() { done(err) }()

	res, err = s.messageRepo.GetByReceiverID(ctx, userID, limit, offset)
	return
}

func (s *MessageService) GetHistoryMessages(ctx context.Context, userID, targetID string, beforeTime int64, limit int) (res []*entity.Message, err error) {
	done := logger.StartStep("MessageService.GetHistoryMessages", "user_id", userID, "target", targetID, "limit", limit)
	defer func() { done(err) }()

	res, err = s.messageRepo.GetHistory(ctx, userID, targetID, beforeTime, limit)
	return
}

func (s *MessageService) MarkAsRead(ctx context.Context, messageID, userID string) (err error) {
	done := logger.StartStep("MessageService.MarkAsRead", "msg_id", messageID, "user", userID)
	defer func() { done(err) }()

	var message *entity.Message
	message, err = s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		err = ErrMessageNotFound
		return
	}

	if message.ReceiverID != userID {
		err = errors.New("not message receiver")
		return
	}

	if message.IsRead {
		return
	}

	if err = s.messageRepo.MarkAsRead(ctx, messageID); err != nil {
		return
	}

	if e := s.messageCache.DecrUnreadCount(ctx, userID); e != nil {
		logger.Warnw("MarkAsRead: decr unread failed", "component", "message_service", "err", e, "user", userID)
	}

	s.eventPublisher.Publish(&event.MessageReadEvent{
		BaseEvent: pkgevent.BaseEvent{
			EventType:   "message.read",
			OccurredAt:  time.Now(),
			AggregateID: messageID,
		},
		MessageID: messageID,
		UserID:    userID,
	})

	logger.Infow("MarkAsRead: marked", "component", "message_service", "msg_id", messageID, "user", userID)
	return
}

func (s *MessageService) MarkAllAsRead(ctx context.Context, userID, targetID string) (err error) {
	done := logger.StartStep("MessageService.MarkAllAsRead", "user", userID, "target", targetID)
	defer func() { done(err) }()

	err = s.messageRepo.MarkAllAsRead(ctx, userID, targetID)
	return
}

func (s *MessageService) GetUnreadCount(ctx context.Context, userID string) (count int64, err error) {
	done := logger.StartStep("MessageService.GetUnreadCount", "user", userID)
	defer func() { done(err) }()

	count, err = s.messageCache.GetUnreadCount(ctx, userID)
	if err != nil {
		return s.messageRepo.GetUnreadCount(ctx, userID)
	}
	return
}

func (s *MessageService) GetOnlineStatus(ctx context.Context, userIDs []string) (result map[string]bool, err error) {
	done := logger.StartStep("MessageService.GetOnlineStatus", "count", len(userIDs))
	defer func() { done(err) }()

	var onlineUsers map[string]bool
	onlineUsers, err = s.messageCache.GetOnlineUsers(ctx)
	if err != nil {
		result = make(map[string]bool)
		return
	}

	result = make(map[string]bool)
	for _, userID := range userIDs {
		result[userID] = onlineUsers[userID]
	}
	return
}

func (s *MessageService) RevokeMessage(ctx context.Context, messageID, userID string) (err error) {
	done := logger.StartStep("MessageService.RevokeMessage", "msg_id", messageID, "user", userID)
	defer func() { done(err) }()

	var message *entity.Message
	message, err = s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		err = ErrMessageNotFound
		return
	}

	if message.SenderID != userID {
		err = ErrNotSender
		return
	}

	if message.IsRevoked {
		err = ErrAlreadyRevoked
		return
	}

	if err = s.messageRepo.Revoke(ctx, messageID); err != nil {
		return
	}

	s.eventPublisher.Publish(&event.MessageRevokedEvent{
		BaseEvent: pkgevent.BaseEvent{
			EventType:   "message.revoked",
			OccurredAt:  time.Now(),
			AggregateID: messageID,
		},
		MessageID: messageID,
		UserID:    userID,
	})

	logger.Infow("RevokeMessage: revoked", "component", "message_service", "msg_id", messageID)
	return
}

func (s *MessageService) SetUserOnline(ctx context.Context, userID string) (err error) {
	done := logger.StartStep("MessageService.SetUserOnline", "user", userID)
	defer func() { done(err) }()

	if err = s.messageCache.SetOnlineUser(ctx, userID); err != nil {
		return
	}

	s.eventPublisher.Publish(&event.UserOnlineEvent{
		BaseEvent: pkgevent.BaseEvent{
			EventType:   "user.online",
			OccurredAt:  time.Now(),
			AggregateID: userID,
		},
		UserID: userID,
	})

	logger.Infow("SetUserOnline: set", "component", "message_service", "user", userID)
	return
}

func (s *MessageService) SetUserOffline(ctx context.Context, userID string) (err error) {
	done := logger.StartStep("MessageService.SetUserOffline", "user", userID)
	defer func() { done(err) }()

	if err = s.messageCache.RemoveOnlineUser(ctx, userID); err != nil {
		return
	}

	s.eventPublisher.Publish(&event.UserOfflineEvent{
		BaseEvent: pkgevent.BaseEvent{
			EventType:   "user.offline",
			OccurredAt:  time.Now(),
			AggregateID: userID,
		},
		UserID: userID,
	})

	logger.Infow("SetUserOffline: removed", "component", "message_service", "user", userID)
	return
}
