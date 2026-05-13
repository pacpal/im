package repository

import (
	"IM/services/group-service/domain/entity"
	"context"
)

type GroupRepository interface {
	Create(ctx context.Context, group *entity.Group) error
	GetByID(ctx context.Context, id string) (*entity.Group, error)
	Update(ctx context.Context, group *entity.Group) error
	Delete(ctx context.Context, id string) error
	GetGroupsByUserID(ctx context.Context, userID string) ([]*entity.Group, error)
}

type GroupMemberRepository interface {
	AddMember(ctx context.Context, gm *entity.GroupMember) error
	RemoveMember(ctx context.Context, groupID, userID string) error
	IsMember(ctx context.Context, groupID, userID string) (bool, error)
	GetMembers(ctx context.Context, groupID string) ([]*entity.User, error)
	GetMemberIDs(ctx context.Context, groupID string) ([]string, error)
	GetRole(ctx context.Context, groupID, userID string) (int16, error)
	UpdateRole(ctx context.Context, groupID, userID string, role int16) error
}

type GroupJoinRequestRepository interface {
	Create(ctx context.Context, req *entity.GroupJoinRequest) error
	GetByID(ctx context.Context, id string) (*entity.GroupJoinRequest, error)
	GetPendingRequests(ctx context.Context, groupID string) ([]*entity.GroupJoinRequest, error)
	UpdateStatus(ctx context.Context, id, status string) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, userID, groupID string) (bool, error)
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*entity.User, error)
	Exists(ctx context.Context, id string) (bool, error)
}

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Tele   string `json:"tele"`
	Status int    `json:"status"`
}