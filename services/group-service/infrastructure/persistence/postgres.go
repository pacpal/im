package persistence

import (
	"IM/services/group-service/domain/entity"
	"IM/services/group-service/infrastructure/repository"
	"IM/server/model"
	"IM/server/repository/postgres"
	"context"
)

type PostgresGroupRepo struct {
	delegate *postgres.GroupRepo
}

func NewPostgresGroupRepo(delegate *postgres.GroupRepo) repository.GroupRepository {
	return &PostgresGroupRepo{delegate: delegate}
}

func (r *PostgresGroupRepo) toDomain(g *model.Group) *entity.Group {
	if g == nil {
		return nil
	}
	return &entity.Group{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		OwnerID:     g.OwnerID,
		Type:        g.Type,
		ImageURL:    g.ImageURL,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}

func (r *PostgresGroupRepo) toDomainList(groups []*model.Group) []*entity.Group {
	result := make([]*entity.Group, len(groups))
	for i, g := range groups {
		result[i] = r.toDomain(g)
	}
	return result
}

func (r *PostgresGroupRepo) toModel(e *entity.Group) *model.Group {
	return &model.Group{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		OwnerID:     e.OwnerID,
		Type:        e.Type,
		ImageURL:    e.ImageURL,
	}
}

func (r *PostgresGroupRepo) Create(ctx context.Context, group *entity.Group) error {
	return r.delegate.Create(ctx, r.toModel(group))
}

func (r *PostgresGroupRepo) GetByID(ctx context.Context, id string) (*entity.Group, error) {
	g, err := r.delegate.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toDomain(g), nil
}

func (r *PostgresGroupRepo) Update(ctx context.Context, group *entity.Group) error {
	return r.delegate.Update(ctx, r.toModel(group))
}

func (r *PostgresGroupRepo) Delete(ctx context.Context, id string) error {
	return r.delegate.Delete(ctx, id)
}

func (r *PostgresGroupRepo) GetGroupsByUserID(ctx context.Context, userID string) ([]*entity.Group, error) {
	groups, err := r.delegate.GetGroupsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return r.toDomainList(groups), nil
}

type PostgresGroupMemberRepo struct {
	delegate *postgres.GroupMemberRepo
}

func NewPostgresGroupMemberRepo(delegate *postgres.GroupMemberRepo) repository.GroupMemberRepository {
	return &PostgresGroupMemberRepo{delegate: delegate}
}

func (r *PostgresGroupMemberRepo) toDomainMemberList(members []*model.User) []*entity.User {
	result := make([]*entity.User, len(members))
	for i, m := range members {
		result[i] = &entity.User{
			ID:     m.ID,
			Name:   m.Name,
			Tele:   m.Tele,
			Status: m.Status,
		}
	}
	return result
}

func (r *PostgresGroupMemberRepo) AddMember(ctx context.Context, gm *entity.GroupMember) error {
	return r.delegate.AddMember(ctx, &model.GroupMember{
		GroupID:  gm.GroupID,
		UserID:   gm.UserID,
		Role:     gm.Role,
		Nickname: gm.Nickname,
	})
}

func (r *PostgresGroupMemberRepo) RemoveMember(ctx context.Context, groupID, userID string) error {
	return r.delegate.RemoveMember(ctx, groupID, userID)
}

func (r *PostgresGroupMemberRepo) IsMember(ctx context.Context, groupID, userID string) (bool, error) {
	return r.delegate.IsMember(ctx, groupID, userID)
}

func (r *PostgresGroupMemberRepo) GetMembers(ctx context.Context, groupID string) ([]*entity.User, error) {
	members, err := r.delegate.GetMembers(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return r.toDomainMemberList(members), nil
}

func (r *PostgresGroupMemberRepo) GetMemberIDs(ctx context.Context, groupID string) ([]string, error) {
	return r.delegate.GetMemberIDs(ctx, groupID)
}

func (r *PostgresGroupMemberRepo) GetRole(ctx context.Context, groupID, userID string) (int16, error) {
	return r.delegate.GetRole(ctx, groupID, userID)
}

func (r *PostgresGroupMemberRepo) UpdateRole(ctx context.Context, groupID, userID string, role int16) error {
	return r.delegate.UpdateRole(ctx, groupID, userID, role)
}

type PostgresGroupJoinRequestRepo struct {
	delegate *postgres.GroupJoinRequestRepo
}

func NewPostgresGroupJoinRequestRepo(delegate *postgres.GroupJoinRequestRepo) repository.GroupJoinRequestRepository {
	return &PostgresGroupJoinRequestRepo{delegate: delegate}
}

func (r *PostgresGroupJoinRequestRepo) toDomain(req *model.GroupJoinRequest) *entity.GroupJoinRequest {
	if req == nil {
		return nil
	}
	return &entity.GroupJoinRequest{
		ID:        req.ID,
		UserID:    req.UserID,
		GroupID:   req.GroupID,
		Reason:    req.Reason,
		Status:    req.Status,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
	}
}

func (r *PostgresGroupJoinRequestRepo) toDomainList(reqs []*model.GroupJoinRequest) []*entity.GroupJoinRequest {
	result := make([]*entity.GroupJoinRequest, len(reqs))
	for i, req := range reqs {
		result[i] = r.toDomain(req)
	}
	return result
}

func (r *PostgresGroupJoinRequestRepo) Create(ctx context.Context, req *entity.GroupJoinRequest) error {
	return r.delegate.Create(ctx, &model.GroupJoinRequest{
		ID:      req.ID,
		UserID:  req.UserID,
		GroupID: req.GroupID,
		Reason:  req.Reason,
		Status:  req.Status,
	})
}

func (r *PostgresGroupJoinRequestRepo) GetByID(ctx context.Context, id string) (*entity.GroupJoinRequest, error) {
	req, err := r.delegate.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toDomain(req), nil
}

func (r *PostgresGroupJoinRequestRepo) GetPendingRequests(ctx context.Context, groupID string) ([]*entity.GroupJoinRequest, error) {
	reqs, err := r.delegate.GetPendingRequests(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return r.toDomainList(reqs), nil
}

func (r *PostgresGroupJoinRequestRepo) UpdateStatus(ctx context.Context, id, status string) error {
	return r.delegate.UpdateStatus(ctx, id, status)
}

func (r *PostgresGroupJoinRequestRepo) Delete(ctx context.Context, id string) error {
	return r.delegate.Delete(ctx, id)
}

func (r *PostgresGroupJoinRequestRepo) Exists(ctx context.Context, userID, groupID string) (bool, error) {
	return r.delegate.Exists(ctx, userID, groupID)
}

type PostgresUserRepo struct {
	delegate *postgres.UserRepo
}

func NewPostgresUserRepo(delegate *postgres.UserRepo) repository.UserRepository {
	return &PostgresUserRepo{delegate: delegate}
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id string) (*entity.User, error) {
	u, err := r.delegate.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.User{
		ID:     u.ID,
		Name:   u.Name,
		Tele:   u.Tele,
		Status: u.Status,
	}, nil
}

func (r *PostgresUserRepo) Exists(ctx context.Context, id string) (bool, error) {
	return r.delegate.Exists(ctx, id)
}