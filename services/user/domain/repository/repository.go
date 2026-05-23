// Package repository 定义 user 域的持久化接口，供上层业务层依赖。
package repository

import (
	"IM/services/user/domain/entity"
	"context"
)

// UserRepository 定义用户持久化所需的方法集合。
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByTele(ctx context.Context, tele string) (*entity.User, error)
	GetByName(ctx context.Context, name string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
	ExistsByTele(ctx context.Context, tele string) (bool, error)
	GetByIDs(ctx context.Context, ids []string) ([]*entity.User, error)
}

// FriendshipRepository 定义好友关系相关的持久化方法。
type FriendshipRepository interface {
	Create(ctx context.Context, f *entity.Friendship) error
	Delete(ctx context.Context, userID, friendID string) error
	Exists(ctx context.Context, userID, friendID string) (bool, error)
	GetFriends(ctx context.Context, userID string) ([]*entity.User, error)
	GetFriendIDs(ctx context.Context, userID string) ([]string, error)
}

// FriendRequestRepository 定义好友申请相关的持久化方法。
type FriendRequestRepository interface {
	Create(ctx context.Context, req *entity.FriendRequest) error
	GetByID(ctx context.Context, id string) (*entity.FriendRequest, error)
	GetPendingRequests(ctx context.Context, uid string) ([]*entity.FriendRequest, error)
	UpdateStatus(ctx context.Context, id string, status entity.FriendRequestStatus) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, fromUID, toUID string) (bool, error)
}
