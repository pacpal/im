// Package repository 定义 group 域的持久化接口。
package repository

import (
	"IM/services/group/domain/entity"
	"context"
)

// GroupRepository 提供群组基本持久化操作接口。
type GroupRepository interface {
	Create(ctx context.Context, group *entity.Group) error
	GetByID(ctx context.Context, id string) (*entity.Group, error)
	Update(ctx context.Context, group *entity.Group) error
	Delete(ctx context.Context, id string) error
	GetByOwnerID(ctx context.Context, ownerID string) ([]*entity.Group, error)
	GetByUserID(ctx context.Context, userID string) ([]*entity.Group, error)
}

// GroupMemberRepository 提供群成员相关的持久化接口。
type GroupMemberRepository interface {
	Create(ctx context.Context, member *entity.GroupMember) error
	GetByGroupID(ctx context.Context, groupID string) ([]*entity.GroupMember, error)
	GetByUserID(ctx context.Context, userID string) ([]*entity.GroupMember, error)
	GetByGroupAndUserID(ctx context.Context, groupID, userID string) (*entity.GroupMember, error)
	Delete(ctx context.Context, groupID, userID string) error
	Exists(ctx context.Context, groupID, userID string) (bool, error)
	UpdateRole(ctx context.Context, groupID, userID string, role entity.MemberRole) error
}

// GroupJoinRequestRepository 定义群组加入请求的持久化接口。
type GroupJoinRequestRepository interface {
	Create(ctx context.Context, req *entity.GroupJoinRequest) error
	GetByID(ctx context.Context, id string) (*entity.GroupJoinRequest, error)
	GetPendingByGroupID(ctx context.Context, groupID string) ([]*entity.GroupJoinRequest, error)
	GetPendingByUserID(ctx context.Context, userID string) ([]*entity.GroupJoinRequest, error)
	UpdateStatus(ctx context.Context, id string, status entity.RequestStatus) error
	Exists(ctx context.Context, userID, groupID string) (bool, error)
}
