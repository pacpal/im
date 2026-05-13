package repository

import (
	"IM/services/message-service/domain/entity"
	"context"
)

type MessageRepository interface {
	Create(ctx context.Context, msg *entity.Message) error
	GetByID(ctx context.Context, id string) (*entity.Message, error)
	GetOfflineMessages(ctx context.Context, userID string, limit, offset int) ([]*entity.Message, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	GetBySender(ctx context.Context, senderID string, limit, offset int) ([]*entity.Message, error)
	GetByReceiver(ctx context.Context, receiverID string, limit, offset int) ([]*entity.Message, error)
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
}

type FriendshipRepository interface {
	Exists(ctx context.Context, userID, friendID string) (bool, error)
}

type GroupMemberRepository interface {
	IsMember(ctx context.Context, groupID, userID string) (bool, error)
	GetMemberIDs(ctx context.Context, groupID string) ([]string, error)
}