package repo

import (
	"IM/server/model"
	"context"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, uid string) (*model.User, error)
	GetByTele(ctx context.Context, tele string) (*model.User, error)
	GetByName(ctx context.Context, name string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, uid string) error
	Exists(ctx context.Context, uid string) (bool, error)
	ExistsByTele(ctx context.Context, tele string) (bool, error)
	GetByIDs(ctx context.Context, uids []string) ([]*model.User, error)
}

type FriendshipRepo interface {
	Create(ctx context.Context, friendship *model.Friendship) error
	Delete(ctx context.Context, userID, friendID string) error
	Exists(ctx context.Context, userID, friendID string) (bool, error)
	GetFriends(ctx context.Context, uid string) ([]*model.User, error)
	GetFriendIDs(ctx context.Context, uid string) ([]string, error)
}

type GroupRepo interface {
	Create(ctx context.Context, group *model.Group) error
	GetByID(ctx context.Context, gid string) (*model.Group, error)
	Update(ctx context.Context, group *model.Group) error
	Delete(ctx context.Context, gid string) error
	GetGroupsByUserID(ctx context.Context, uid string) ([]*model.Group, error)
}

type GroupMemberRepo interface {
	AddMember(ctx context.Context, gm *model.GroupMember) error
	RemoveMember(ctx context.Context, gid, uid string) error
	IsMember(ctx context.Context, gid, uid string) (bool, error)
	GetMembers(ctx context.Context, gid string) ([]*model.User, error)
	GetMemberIDs(ctx context.Context, gid string) ([]string, error)
	GetRole(ctx context.Context, gid, uid string) (int16, error)
	UpdateRole(ctx context.Context, gid, uid string, role int16) error
}

type MessageRepo interface {
	Create(ctx context.Context, msg *model.Message) error
	GetByID(ctx context.Context, msgID string) (*model.Message, error)
	GetOfflineMessages(ctx context.Context, uid string, limit, offset int) ([]*model.Message, error)
	MarkAsRead(ctx context.Context, msgID string) error
	MarkAllAsRead(ctx context.Context, uid string) error
	GetBySender(ctx context.Context, senderID string, limit, offset int) ([]*model.Message, error)
	GetByReceiver(ctx context.Context, receiverID string, limit, offset int) ([]*model.Message, error)
	GetUnreadCount(ctx context.Context, uid string) (int64, error)
}

type FriendRequestRepo interface {
	Create(ctx context.Context, req *model.FriendRequest) error
	GetByID(ctx context.Context, reqID string) (*model.FriendRequest, error)
	GetPendingRequests(ctx context.Context, uid string) ([]*model.FriendRequest, error)
	UpdateStatus(ctx context.Context, reqID, status string) error
	Delete(ctx context.Context, reqID string) error
	Exists(ctx context.Context, fromUID, toUID string) (bool, error)
}

type GroupJoinRequestRepo interface {
	Create(ctx context.Context, req *model.GroupJoinRequest) error
	GetByID(ctx context.Context, reqID string) (*model.GroupJoinRequest, error)
	GetPendingRequests(ctx context.Context, gid string) ([]*model.GroupJoinRequest, error)
	UpdateStatus(ctx context.Context, reqID, status string) error
	Delete(ctx context.Context, reqID string) error
	Exists(ctx context.Context, userID, groupID string) (bool, error)
}
