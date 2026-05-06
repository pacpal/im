// Package repo is for store
package repo

import (
	"IM/server/model"
	"context"
)

type UserRepo interface {
	GetByID(ctx context.Context, uid string) (bool, error)
	GetByName(ctx context.Context, name string) (bool, error)
	CreateNew(ctx context.Context, uid, name string, password []byte) (bool, error)
	GetUserByID(ctx context.Context, uid string) (*model.User, error)
	GetUserByTele(ctx context.Context, tele string) (*model.User, error)
	SaveUser(ctx context.Context, user *model.User) error
	//...
}
type GroupRepo interface {
	GetGroupByID(ctx context.Context, gid string) (*model.Group, error)
	IsMember(ctx context.Context, gid, uid string) (bool, error)
	SaveGroup(ctx context.Context, user *model.Group) error
}
