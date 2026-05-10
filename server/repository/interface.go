// Package repo is for store
package repo

import (
	"IM/server/model"
	"context"
)

type UserRepo interface {
	Create(ctx context.Context, uid, name string, password []byte) (bool, error)
	GetUserByName(ctx context.Context, name string) (*model.User, error)
	GetUserByID(ctx context.Context, uid string) (*model.User, error)
	GetUserByTele(ctx context.Context, tele string) (*model.User, error)
	RefreshUser(ctx context.Context, user *model.User) error
	//...
}
type GroupRepo interface {
	GetGroupByID(ctx context.Context, gid string) (*model.Group, error)
	IsMember(ctx context.Context, gid, uid string) (bool, error)
	SaveGroup(ctx context.Context, user *model.Group) error
}
type MsgRepo interface {
	GetOfflineMsgs(ctx context.Context, uid string) (*[]model.Message, error)
	ClearOfflineMsgs(ctx context.Context, uid string)
	GetGroupByID(ctx context.Context, gid string) (*model.Group, error)
}
