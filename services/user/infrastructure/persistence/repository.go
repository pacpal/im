package persistence

import (
	"IM/services/user/domain/entity"
	"IM/services/user/domain/repository"
	"IM/services/user/infrastructure/persistence/model"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	m := toUserModel(user)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return toUserEntity(&m), nil
}

func (r *UserRepository) GetByTele(ctx context.Context, tele string) (*entity.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).Where("tele = ?", tele).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return toUserEntity(&m), nil
}

func (r *UserRepository) GetByName(ctx context.Context, name string) (*entity.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return toUserEntity(&m), nil
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	m := toUserModel(user)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}

func (r *UserRepository) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) ExistsByTele(ctx context.Context, tele string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("tele = ?", tele).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) GetByIDs(ctx context.Context, ids []string) ([]*entity.User, error) {
	var models []*model.User
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&models).Error; err != nil {
		return nil, err
	}
	users := make([]*entity.User, len(models))
	for i, m := range models {
		users[i] = toUserEntity(m)
	}
	return users, nil
}

type FriendshipRepository struct {
	db *gorm.DB
}

func NewFriendshipRepository(db *gorm.DB) repository.FriendshipRepository {
	return &FriendshipRepository{db: db}
}

func (r *FriendshipRepository) Create(ctx context.Context, f *entity.Friendship) error {
	m := &model.Friendship{
		UserID:    f.UserID,
		FriendID:  f.FriendID,
		Status:    int(f.Status),
		CreatedAt: f.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *FriendshipRepository) Delete(ctx context.Context, userID, friendID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Delete(&model.Friendship{}).Error
}

func (r *FriendshipRepository) Exists(ctx context.Context, userID, friendID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Friendship{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *FriendshipRepository) GetFriends(ctx context.Context, userID string) ([]*entity.User, error) {
	var friendIDs []string
	if err := r.db.WithContext(ctx).Model(&model.Friendship{}).
		Where("user_id = ? AND status = ?", userID, entity.FriendshipStatusActive).
		Pluck("friend_id", &friendIDs).Error; err != nil {
		return nil, err
	}

	if len(friendIDs) == 0 {
		return []*entity.User{}, nil
	}

	var users []*model.User
	if err := r.db.WithContext(ctx).Where("id IN ?", friendIDs).Find(&users).Error; err != nil {
		return nil, err
	}

	result := make([]*entity.User, len(users))
	for i, u := range users {
		result[i] = toUserEntity(u)
	}
	return result, nil
}

func (r *FriendshipRepository) GetFriendIDs(ctx context.Context, userID string) ([]string, error) {
	var friendIDs []string
	if err := r.db.WithContext(ctx).Model(&model.Friendship{}).
		Where("user_id = ? AND status = ?", userID, entity.FriendshipStatusActive).
		Pluck("friend_id", &friendIDs).Error; err != nil {
		return nil, err
	}
	return friendIDs, nil
}

type FriendRequestRepository struct {
	db *gorm.DB
}

func NewFriendRequestRepository(db *gorm.DB) repository.FriendRequestRepository {
	return &FriendRequestRepository{db: db}
}

func (r *FriendRequestRepository) Create(ctx context.Context, req *entity.FriendRequest) error {
	m := &model.FriendRequest{
		ID:        req.ID,
		FromUID:   req.FromUID,
		ToUID:     req.ToUID,
		Reason:    req.Reason,
		Status:    string(req.Status),
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *FriendRequestRepository) GetByID(ctx context.Context, id string) (*entity.FriendRequest, error) {
	var m model.FriendRequest
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("request not found")
		}
		return nil, err
	}
	return toFriendRequestEntity(&m), nil
}

func (r *FriendRequestRepository) GetPendingRequests(ctx context.Context, uid string) ([]*entity.FriendRequest, error) {
	var models []*model.FriendRequest
	if err := r.db.WithContext(ctx).
		Where("to_uid = ? AND status = ?", uid, entity.FriendRequestStatusPending).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]*entity.FriendRequest, len(models))
	for i, m := range models {
		result[i] = toFriendRequestEntity(m)
	}
	return result, nil
}

func (r *FriendRequestRepository) UpdateStatus(ctx context.Context, id string, status entity.FriendRequestStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.FriendRequest{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     string(status),
			"updated_at": time.Now(),
		}).Error
}

func (r *FriendRequestRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.FriendRequest{}, "id = ?", id).Error
}

func (r *FriendRequestRepository) Exists(ctx context.Context, fromUID, toUID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.FriendRequest{}).
		Where("from_uid = ? AND to_uid = ? AND status = ?", fromUID, toUID, entity.FriendRequestStatusPending).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func toUserEntity(m *model.User) *entity.User {
	return &entity.User{
		ID:        m.ID,
		Name:      m.Name,
		Tele:      m.Tele,
		Password:  m.Password,
		AvatarURL: m.AvatarURL,
		Status:    entity.UserStatus(m.Status),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func toUserModel(e *entity.User) *model.User {
	return &model.User{
		ID:        e.ID,
		Name:      e.Name,
		Tele:      e.Tele,
		Password:  e.Password,
		AvatarURL: e.AvatarURL,
		Status:    int(e.Status),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func toFriendRequestEntity(m *model.FriendRequest) *entity.FriendRequest {
	return &entity.FriendRequest{
		ID:        m.ID,
		FromUID:   m.FromUID,
		ToUID:     m.ToUID,
		Reason:    m.Reason,
		Status:    entity.FriendRequestStatus(m.Status),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
