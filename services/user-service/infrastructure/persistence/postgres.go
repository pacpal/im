package persistence

import (
	"IM/services/user-service/domain/entity"
	repo "IM/services/user-service/infrastructure/repository"
	"IM/server/model"
	"IM/server/repository/postgres"
	"context"
)

type PostgresUserRepo struct {
	delegate *postgres.UserRepo
}

func NewPostgresUserRepo(delegate *postgres.UserRepo) repo.UserRepository {
	return &PostgresUserRepo{delegate: delegate}
}

func (r *PostgresUserRepo) toDomain(u *model.User) *entity.User {
	if u == nil {
		return nil
	}
	return &entity.User{
		ID:        u.ID,
		Name:      u.Name,
		Tele:      u.Tele,
		Password:  u.Password,
		AvatarURL: u.AvatarURL,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (r *PostgresUserRepo) toDomainList(users []*model.User) []*entity.User {
	result := make([]*entity.User, len(users))
	for i, u := range users {
		result[i] = r.toDomain(u)
	}
	return result
}

func (r *PostgresUserRepo) toModel(e *entity.User) *model.User {
	return &model.User{
		ID:        e.ID,
		Name:      e.Name,
		Tele:      e.Tele,
		Password:  e.Password,
		AvatarURL: e.AvatarURL,
		Status:    e.Status,
	}
}

func (r *PostgresUserRepo) Create(ctx context.Context, user *entity.User) error {
	return r.delegate.Create(ctx, r.toModel(user))
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id string) (*entity.User, error) {
	u, err := r.delegate.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toDomain(u), nil
}

func (r *PostgresUserRepo) GetByTele(ctx context.Context, tele string) (*entity.User, error) {
	u, err := r.delegate.GetByTele(ctx, tele)
	if err != nil {
		return nil, err
	}
	return r.toDomain(u), nil
}

func (r *PostgresUserRepo) GetByName(ctx context.Context, name string) (*entity.User, error) {
	u, err := r.delegate.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return r.toDomain(u), nil
}

func (r *PostgresUserRepo) Update(ctx context.Context, user *entity.User) error {
	return r.delegate.Update(ctx, r.toModel(user))
}

func (r *PostgresUserRepo) Delete(ctx context.Context, id string) error {
	return r.delegate.Delete(ctx, id)
}

func (r *PostgresUserRepo) Exists(ctx context.Context, id string) (bool, error) {
	return r.delegate.Exists(ctx, id)
}

func (r *PostgresUserRepo) ExistsByTele(ctx context.Context, tele string) (bool, error) {
	return r.delegate.ExistsByTele(ctx, tele)
}

func (r *PostgresUserRepo) GetByIDs(ctx context.Context, ids []string) ([]*entity.User, error) {
	users, err := r.delegate.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return r.toDomainList(users), nil
}

type PostgresFriendshipRepo struct {
	delegate *postgres.FriendshipRepo
}

func NewPostgresFriendshipRepo(delegate *postgres.FriendshipRepo) repo.FriendshipRepository {
	return &PostgresFriendshipRepo{delegate: delegate}
}

func (r *PostgresFriendshipRepo) toDomainFriendship(f *model.Friendship) *entity.Friendship {
	if f == nil {
		return nil
	}
	return &entity.Friendship{
		UserID:    f.UserID,
		FriendID:  f.FriendID,
		Status:    f.Status,
		CreatedAt: f.CreatedAt,
	}
}

func (r *PostgresFriendshipRepo) toDomainFriendshipList(friends []*model.User) []*entity.User {
	result := make([]*entity.User, len(friends))
	for i, f := range friends {
		result[i] = &entity.User{
			ID:     f.ID,
			Name:   f.Name,
			Tele:   f.Tele,
			Status: f.Status,
		}
	}
	return result
}

func (r *PostgresFriendshipRepo) Create(ctx context.Context, f *entity.Friendship) error {
	return r.delegate.Create(ctx, &model.Friendship{
		UserID:   f.UserID,
		FriendID: f.FriendID,
		Status:   f.Status,
	})
}

func (r *PostgresFriendshipRepo) Delete(ctx context.Context, userID, friendID string) error {
	return r.delegate.Delete(ctx, userID, friendID)
}

func (r *PostgresFriendshipRepo) Exists(ctx context.Context, userID, friendID string) (bool, error) {
	return r.delegate.Exists(ctx, userID, friendID)
}

func (r *PostgresFriendshipRepo) GetFriends(ctx context.Context, userID string) ([]*entity.User, error) {
	users, err := r.delegate.GetFriends(ctx, userID)
	if err != nil {
		return nil, err
	}
	return r.toDomainFriendshipList(users), nil
}

func (r *PostgresFriendshipRepo) GetFriendIDs(ctx context.Context, userID string) ([]string, error) {
	return r.delegate.GetFriendIDs(ctx, userID)
}

type PostgresFriendRequestRepo struct {
	delegate *postgres.FriendRequestRepo
}

func NewPostgresFriendRequestRepo(delegate *postgres.FriendRequestRepo) repo.FriendRequestRepository {
	return &PostgresFriendRequestRepo{delegate: delegate}
}

func (r *PostgresFriendRequestRepo) toDomain(req *model.FriendRequest) *entity.FriendRequest {
	if req == nil {
		return nil
	}
	return &entity.FriendRequest{
		ID:        req.ID,
		FromUID:   req.FromUID,
		ToUID:     req.ToUID,
		Reason:    req.Reason,
		Status:    req.Status,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
	}
}

func (r *PostgresFriendRequestRepo) toDomainList(reqs []*model.FriendRequest) []*entity.FriendRequest {
	result := make([]*entity.FriendRequest, len(reqs))
	for i, req := range reqs {
		result[i] = r.toDomain(req)
	}
	return result
}

func (r *PostgresFriendRequestRepo) Create(ctx context.Context, req *entity.FriendRequest) error {
	return r.delegate.Create(ctx, &model.FriendRequest{
		ID:      req.ID,
		FromUID: req.FromUID,
		ToUID:   req.ToUID,
		Reason:  req.Reason,
		Status:  req.Status,
	})
}

func (r *PostgresFriendRequestRepo) GetByID(ctx context.Context, id string) (*entity.FriendRequest, error) {
	req, err := r.delegate.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toDomain(req), nil
}

func (r *PostgresFriendRequestRepo) GetPendingRequests(ctx context.Context, uid string) ([]*entity.FriendRequest, error) {
	reqs, err := r.delegate.GetPendingRequests(ctx, uid)
	if err != nil {
		return nil, err
	}
	return r.toDomainList(reqs), nil
}

func (r *PostgresFriendRequestRepo) UpdateStatus(ctx context.Context, id, status string) error {
	return r.delegate.UpdateStatus(ctx, id, status)
}

func (r *PostgresFriendRequestRepo) Delete(ctx context.Context, id string) error {
	return r.delegate.Delete(ctx, id)
}

func (r *PostgresFriendRequestRepo) Exists(ctx context.Context, fromUID, toUID string) (bool, error) {
	return r.delegate.Exists(ctx, fromUID, toUID)
}