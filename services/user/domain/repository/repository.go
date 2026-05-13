package repository

import (
	"IM/services/user/domain/entity"
	"context"
)

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

type FriendshipRepository interface {
	Create(ctx context.Context, f *entity.Friendship) error
	Delete(ctx context.Context, userID, friendID string) error
	Exists(ctx context.Context, userID, friendID string) (bool, error)
	GetFriends(ctx context.Context, userID string) ([]*entity.User, error)
	GetFriendIDs(ctx context.Context, userID string) ([]string, error)
}

type FriendRequestRepository interface {
	Create(ctx context.Context, req *entity.FriendRequest) error
	GetByID(ctx context.Context, id string) (*entity.FriendRequest, error)
	GetPendingRequests(ctx context.Context, uid string) ([]*entity.FriendRequest, error)
	UpdateStatus(ctx context.Context, id string, status entity.FriendRequestStatus) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, fromUID, toUID string) (bool, error)
}
