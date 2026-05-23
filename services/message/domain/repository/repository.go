// Package repository 定义消息服务使用的仓储与缓存接口。
package repository

import (
	"IM/services/message/domain/entity"
	"context"
)

// MessageRepository 抽象消息的持久化操作。
type MessageRepository interface {
	Create(ctx context.Context, message *entity.Message) error
	GetByID(ctx context.Context, id string) (*entity.Message, error)
	GetByReceiverID(ctx context.Context, receiverID string, limit, offset int) ([]*entity.Message, error)
	GetHistory(ctx context.Context, userID, targetID string, beforeTime int64, limit int) ([]*entity.Message, error)
	GetUnreadByReceiverID(ctx context.Context, receiverID string) ([]*entity.Message, error)
	GetUnreadCount(ctx context.Context, receiverID string) (int64, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, receiverID, senderID string) error
	Revoke(ctx context.Context, id string) error
}

// MessageCache 抽象缓存层（Redis）的接口，用于在线用户和未读计数管理。
type MessageCache interface {
	GetOnlineUsers(ctx context.Context) (map[string]bool, error)
	SetOnlineUser(ctx context.Context, userID string) error
	RemoveOnlineUser(ctx context.Context, userID string) error
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
	IncrUnreadCount(ctx context.Context, userID string) error
	DecrUnreadCount(ctx context.Context, userID string) error
}
