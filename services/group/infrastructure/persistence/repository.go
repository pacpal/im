package persistence

import (
	"IM/services/group/domain/entity"
	"IM/services/group/domain/repository"
	"IM/services/group/infrastructure/persistence/model"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) repository.GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(ctx context.Context, group *entity.Group) error {
	m := toGroupModel(group)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *GroupRepository) GetByID(ctx context.Context, id string) (*entity.Group, error) {
	var m model.Group
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("group not found")
		}
		return nil, err
	}
	return toGroupEntity(&m), nil
}

func (r *GroupRepository) Update(ctx context.Context, group *entity.Group) error {
	m := toGroupModel(group)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *GroupRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Group{}, "id = ?", id).Error
}

func (r *GroupRepository) GetByOwnerID(ctx context.Context, ownerID string) ([]*entity.Group, error) {
	var models []*model.Group
	if err := r.db.WithContext(ctx).Where("owner_id = ?", ownerID).Find(&models).Error; err != nil {
		return nil, err
	}
	groups := make([]*entity.Group, len(models))
	for i, m := range models {
		groups[i] = toGroupEntity(m)
	}
	return groups, nil
}

func (r *GroupRepository) GetByUserID(ctx context.Context, userID string) ([]*entity.Group, error) {
	var groupIDs []string
	if err := r.db.WithContext(ctx).Model(&model.GroupMember{}).
		Where("user_id = ?", userID).
		Pluck("group_id", &groupIDs).Error; err != nil {
		return nil, err
	}

	if len(groupIDs) == 0 {
		return []*entity.Group{}, nil
	}

	var models []*model.Group
	if err := r.db.WithContext(ctx).Where("id IN ?", groupIDs).Find(&models).Error; err != nil {
		return nil, err
	}

	groups := make([]*entity.Group, len(models))
	for i, m := range models {
		groups[i] = toGroupEntity(m)
	}
	return groups, nil
}

type GroupMemberRepository struct {
	db *gorm.DB
}

func NewGroupMemberRepository(db *gorm.DB) repository.GroupMemberRepository {
	return &GroupMemberRepository{db: db}
}

func (r *GroupMemberRepository) Create(ctx context.Context, member *entity.GroupMember) error {
	m := &model.GroupMember{
		GroupID:  member.GroupID,
		UserID:   member.UserID,
		Role:     int(member.Role),
		Nickname: member.Nickname,
		JoinedAt: member.JoinedAt,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *GroupMemberRepository) GetByGroupID(ctx context.Context, groupID string) ([]*entity.GroupMember, error) {
	var models []*model.GroupMember
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&models).Error; err != nil {
		return nil, err
	}
	members := make([]*entity.GroupMember, len(models))
	for i, m := range models {
		members[i] = toGroupMemberEntity(m)
	}
	return members, nil
}

func (r *GroupMemberRepository) GetByUserID(ctx context.Context, userID string) ([]*entity.GroupMember, error) {
	var models []*model.GroupMember
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}
	members := make([]*entity.GroupMember, len(models))
	for i, m := range models {
		members[i] = toGroupMemberEntity(m)
	}
	return members, nil
}

func (r *GroupMemberRepository) GetByGroupAndUserID(ctx context.Context, groupID, userID string) (*entity.GroupMember, error) {
	var m model.GroupMember
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("member not found")
		}
		return nil, err
	}
	return toGroupMemberEntity(&m), nil
}

func (r *GroupMemberRepository) Delete(ctx context.Context, groupID, userID string) error {
	return r.db.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Delete(&model.GroupMember{}).Error
}

func (r *GroupMemberRepository) Exists(ctx context.Context, groupID, userID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GroupMemberRepository) UpdateRole(ctx context.Context, groupID, userID string, role entity.MemberRole) error {
	return r.db.WithContext(ctx).
		Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("role", int(role)).Error
}

type GroupJoinRequestRepository struct {
	db *gorm.DB
}

func NewGroupJoinRequestRepository(db *gorm.DB) repository.GroupJoinRequestRepository {
	return &GroupJoinRequestRepository{db: db}
}

func (r *GroupJoinRequestRepository) Create(ctx context.Context, req *entity.GroupJoinRequest) error {
	m := &model.GroupJoinRequest{
		ID:        req.ID,
		UserID:    req.UserID,
		GroupID:   req.GroupID,
		Reason:    req.Reason,
		Status:    string(req.Status),
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *GroupJoinRequestRepository) GetByID(ctx context.Context, id string) (*entity.GroupJoinRequest, error) {
	var m model.GroupJoinRequest
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("request not found")
		}
		return nil, err
	}
	return toGroupJoinRequestEntity(&m), nil
}

func (r *GroupJoinRequestRepository) GetPendingByGroupID(ctx context.Context, groupID string) ([]*entity.GroupJoinRequest, error) {
	var models []*model.GroupJoinRequest
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND status = ?", groupID, entity.RequestStatusPending).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	requests := make([]*entity.GroupJoinRequest, len(models))
	for i, m := range models {
		requests[i] = toGroupJoinRequestEntity(m)
	}
	return requests, nil
}

func (r *GroupJoinRequestRepository) GetPendingByUserID(ctx context.Context, userID string) ([]*entity.GroupJoinRequest, error) {
	var models []*model.GroupJoinRequest
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, entity.RequestStatusPending).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	requests := make([]*entity.GroupJoinRequest, len(models))
	for i, m := range models {
		requests[i] = toGroupJoinRequestEntity(m)
	}
	return requests, nil
}

func (r *GroupJoinRequestRepository) UpdateStatus(ctx context.Context, id string, status entity.RequestStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.GroupJoinRequest{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     string(status),
			"updated_at": time.Now(),
		}).Error
}

func (r *GroupJoinRequestRepository) Exists(ctx context.Context, userID, groupID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.GroupJoinRequest{}).
		Where("user_id = ? AND group_id = ? AND status = ?", userID, groupID, entity.RequestStatusPending).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func toGroupEntity(m *model.Group) *entity.Group {
	return &entity.Group{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		OwnerID:     m.OwnerID,
		Type:        entity.GroupType(m.Type),
		ImageURL:    m.ImageURL,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toGroupModel(e *entity.Group) *model.Group {
	return &model.Group{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		OwnerID:     e.OwnerID,
		Type:        string(e.Type),
		ImageURL:    e.ImageURL,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toGroupMemberEntity(m *model.GroupMember) *entity.GroupMember {
	return &entity.GroupMember{
		GroupID:  m.GroupID,
		UserID:   m.UserID,
		Role:     entity.MemberRole(m.Role),
		Nickname: m.Nickname,
		JoinedAt: m.JoinedAt,
	}
}

func toGroupJoinRequestEntity(m *model.GroupJoinRequest) *entity.GroupJoinRequest {
	return &entity.GroupJoinRequest{
		ID:        m.ID,
		UserID:    m.UserID,
		GroupID:   m.GroupID,
		Reason:    m.Reason,
		Status:    entity.RequestStatus(m.Status),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
