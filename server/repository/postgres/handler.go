package postgres

import (
	"IM/server/model"
	repo "IM/server/repository"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var _ repo.UserRepo = (*UserRepo)(nil)
var _ repo.FriendshipRepo = (*FriendshipRepo)(nil)
var _ repo.GroupRepo = (*GroupRepo)(nil)
var _ repo.GroupMemberRepo = (*GroupMemberRepo)(nil)
var _ repo.MessageRepo = (*MessageRepo)(nil)
var _ repo.FriendRequestRepo = (*FriendRequestRepo)(nil)
var _ repo.GroupJoinRequestRepo = (*GroupJoinRequestRepo)(nil)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) GetByID(ctx context.Context, uid string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", uid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("query user by id failed: %w", err)
	}
	return &user, nil
}

func (r *UserRepo) GetByTele(ctx context.Context, tele string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("tele = ?", tele).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("query user by tele failed: %w", err)
	}
	return &user, nil
}

func (r *UserRepo) GetByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("query user by name failed: %w", err)
	}
	return &user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("update user failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, uid string) error {
	result := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", uid)
	if result.Error != nil {
		return fmt.Errorf("delete user failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

func (r *UserRepo) Exists(ctx context.Context, uid string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", uid).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check user exists failed: %w", err)
	}
	return count > 0, nil
}

func (r *UserRepo) ExistsByTele(ctx context.Context, tele string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("tele = ?", tele).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check user tele exists failed: %w", err)
	}
	return count > 0, nil
}

func (r *UserRepo) GetByIDs(ctx context.Context, uids []string) ([]*model.User, error) {
	var users []*model.User
	err := r.db.WithContext(ctx).Where("id IN ?", uids).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("batch query users failed: %w", err)
	}
	return users, nil
}

type FriendshipRepo struct {
	db *gorm.DB
}

func NewFriendshipRepo(db *gorm.DB) *FriendshipRepo {
	return &FriendshipRepo{db: db}
}

func (r *FriendshipRepo) Create(ctx context.Context, friendship *model.Friendship) error {
	return r.db.WithContext(ctx).Create(friendship).Error
}

func (r *FriendshipRepo) Delete(ctx context.Context, userID, friendID string) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Or("user_id = ? AND friend_id = ?", friendID, userID).
		Delete(&model.Friendship{})
	if result.Error != nil {
		return fmt.Errorf("delete friendship failed: %w", result.Error)
	}
	return nil
}

func (r *FriendshipRepo) Exists(ctx context.Context, userID, friendID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Friendship{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Or("user_id = ? AND friend_id = ?", friendID, userID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check friendship exists failed: %w", err)
	}
	return count > 0, nil
}

func (r *FriendshipRepo) GetFriends(ctx context.Context, uid string) ([]*model.User, error) {
	var users []*model.User
	err := r.db.WithContext(ctx).
		Table("users").
		Joins("JOIN friendships ON users.id = friendships.friend_id").
		Where("friendships.user_id = ? AND friendships.status = ?", uid, 1).
		Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("query friends failed: %w", err)
	}
	return users, nil
}

func (r *FriendshipRepo) GetFriendIDs(ctx context.Context, uid string) ([]string, error) {
	var ids []string
	err := r.db.WithContext(ctx).
		Model(&model.Friendship{}).
		Where("user_id = ? AND status = ?", uid, 1).
		Pluck("friend_id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("query friend ids failed: %w", err)
	}
	return ids, nil
}

type GroupRepo struct {
	db *gorm.DB
}

func NewGroupRepo(db *gorm.DB) *GroupRepo {
	return &GroupRepo{db: db}
}

func (r *GroupRepo) Create(ctx context.Context, group *model.Group) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *GroupRepo) GetByID(ctx context.Context, gid string) (*model.Group, error) {
	var group model.Group
	err := r.db.WithContext(ctx).Where("id = ?", gid).First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrGroupNotFound
		}
		return nil, fmt.Errorf("query group by id failed: %w", err)
	}
	return &group, nil
}

func (r *GroupRepo) Update(ctx context.Context, group *model.Group) error {
	result := r.db.WithContext(ctx).Save(group)
	if result.Error != nil {
		return fmt.Errorf("update group failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrGroupNotFound
	}
	return nil
}

func (r *GroupRepo) Delete(ctx context.Context, gid string) error {
	result := r.db.WithContext(ctx).Delete(&model.Group{}, "id = ?", gid)
	if result.Error != nil {
		return fmt.Errorf("delete group failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrGroupNotFound
	}
	return nil
}

func (r *GroupRepo) GetGroupsByUserID(ctx context.Context, uid string) ([]*model.Group, error) {
	var groups []*model.Group
	err := r.db.WithContext(ctx).
		Table("groups").
		Joins("LEFT JOIN group_members ON groups.id = group_members.group_id").
		Where("groups.owner_id = ? OR group_members.user_id = ?", uid, uid).
		Distinct().
		Find(&groups).Error
	if err != nil {
		return nil, fmt.Errorf("query groups by user id failed: %w", err)
	}
	return groups, nil
}

type GroupMemberRepo struct {
	db *gorm.DB
}

func NewGroupMemberRepo(db *gorm.DB) *GroupMemberRepo {
	return &GroupMemberRepo{db: db}
}

func (r *GroupMemberRepo) AddMember(ctx context.Context, gm *model.GroupMember) error {
	return r.db.WithContext(ctx).Create(gm).Error
}

func (r *GroupMemberRepo) RemoveMember(ctx context.Context, gid, uid string) error {
	result := r.db.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", gid, uid).
		Delete(&model.GroupMember{})
	if result.Error != nil {
		return fmt.Errorf("remove group member failed: %w", result.Error)
	}
	return nil
}

func (r *GroupMemberRepo) IsMember(ctx context.Context, gid, uid string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", gid, uid).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check group member failed: %w", err)
	}
	return count > 0, nil
}

func (r *GroupMemberRepo) GetMembers(ctx context.Context, gid string) ([]*model.User, error) {
	var users []*model.User
	err := r.db.WithContext(ctx).
		Table("users").
		Joins("JOIN group_members ON users.id = group_members.user_id").
		Where("group_members.group_id = ?", gid).
		Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("query group members failed: %w", err)
	}
	return users, nil
}

func (r *GroupMemberRepo) GetMemberIDs(ctx context.Context, gid string) ([]string, error) {
	var ids []string
	err := r.db.WithContext(ctx).
		Model(&model.GroupMember{}).
		Where("group_id = ?", gid).
		Pluck("user_id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("query group member ids failed: %w", err)
	}
	return ids, nil
}

func (r *GroupMemberRepo) GetRole(ctx context.Context, gid, uid string) (int16, error) {
	var role int16
	err := r.db.WithContext(ctx).
		Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", gid, uid).
		Pluck("role", &role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, model.ErrNotMember
		}
		return 0, fmt.Errorf("get member role failed: %w", err)
	}
	return role, nil
}

func (r *GroupMemberRepo) UpdateRole(ctx context.Context, gid, uid string, role int16) error {
	result := r.db.WithContext(ctx).
		Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", gid, uid).
		Update("role", role)
	if result.Error != nil {
		return fmt.Errorf("update member role failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrNotMember
	}
	return nil
}

type MessageRepo struct {
	db *gorm.DB
}

func NewMessageRepo(db *gorm.DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (r *MessageRepo) Create(ctx context.Context, msg *model.Message) error {
	if msg.ID == "" {
		msg.ID = fmt.Sprintf("msg_%d", msg.Time)
	}
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *MessageRepo) GetByID(ctx context.Context, msgID string) (*model.Message, error) {
	var msg model.Message
	err := r.db.WithContext(ctx).Where("id = ?", msgID).First(&msg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrMessageNotFound
		}
		return nil, fmt.Errorf("query message by id failed: %w", err)
	}
	return &msg, nil
}

func (r *MessageRepo) GetOfflineMessages(ctx context.Context, uid string, limit, offset int) ([]*model.Message, error) {
	var messages []*model.Message
	query := r.db.WithContext(ctx).
		Where("receive_id = ? AND is_read = ?", uid, false).
		Order("time ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("query offline messages failed: %w", err)
	}
	return messages, nil
}

func (r *MessageRepo) MarkAsRead(ctx context.Context, msgID string) error {
	return r.db.WithContext(ctx).
		Model(&model.Message{}).
		Where("id = ?", msgID).
		Update("is_read", true).Error
}

func (r *MessageRepo) MarkAllAsRead(ctx context.Context, uid string) error {
	return r.db.WithContext(ctx).
		Model(&model.Message{}).
		Where("receive_id = ? AND is_read = ?", uid, false).
		Update("is_read", true).Error
}

func (r *MessageRepo) GetBySender(ctx context.Context, senderID string, limit, offset int) ([]*model.Message, error) {
	var messages []*model.Message
	query := r.db.WithContext(ctx).
		Where("send_id = ?", senderID).
		Order("time DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("query messages by sender failed: %w", err)
	}
	return messages, nil
}

func (r *MessageRepo) GetByReceiver(ctx context.Context, receiverID string, limit, offset int) ([]*model.Message, error) {
	var messages []*model.Message
	query := r.db.WithContext(ctx).
		Where("receive_id = ?", receiverID).
		Order("time DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("query messages by receiver failed: %w", err)
	}
	return messages, nil
}

func (r *MessageRepo) GetUnreadCount(ctx context.Context, uid string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Message{}).
		Where("receive_id = ? AND is_read = ?", uid, false).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("query unread count failed: %w", err)
	}
	return count, nil
}

type FriendRequestRepo struct {
	db *gorm.DB
}

func NewFriendRequestRepo(db *gorm.DB) *FriendRequestRepo {
	return &FriendRequestRepo{db: db}
}

func (r *FriendRequestRepo) Create(ctx context.Context, req *model.FriendRequest) error {
	return r.db.WithContext(ctx).Create(req).Error
}

func (r *FriendRequestRepo) GetByID(ctx context.Context, reqID string) (*model.FriendRequest, error) {
	var req model.FriendRequest
	err := r.db.WithContext(ctx).Where("id = ?", reqID).First(&req).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrRequestNotFound
		}
		return nil, fmt.Errorf("query friend request failed: %w", err)
	}
	return &req, nil
}

func (r *FriendRequestRepo) GetPendingRequests(ctx context.Context, uid string) ([]*model.FriendRequest, error) {
	var requests []*model.FriendRequest
	err := r.db.WithContext(ctx).
		Where("to_uid = ? AND status = ?", uid, model.FriendRequestPending).
		Order("created_at DESC").
		Find(&requests).Error
	if err != nil {
		return nil, fmt.Errorf("query pending friend requests failed: %w", err)
	}
	return requests, nil
}

func (r *FriendRequestRepo) UpdateStatus(ctx context.Context, reqID, status string) error {
	result := r.db.WithContext(ctx).
		Model(&model.FriendRequest{}).
		Where("id = ?", reqID).
		Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("update friend request status failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrRequestNotFound
	}
	return nil
}

func (r *FriendRequestRepo) Delete(ctx context.Context, reqID string) error {
	result := r.db.WithContext(ctx).Delete(&model.FriendRequest{}, "id = ?", reqID)
	if result.Error != nil {
		return fmt.Errorf("delete friend request failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrRequestNotFound
	}
	return nil
}

func (r *FriendRequestRepo) Exists(ctx context.Context, fromUID, toUID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.FriendRequest{}).
		Where("from_uid = ? AND to_uid = ? AND status = ?", fromUID, toUID, model.FriendRequestPending).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check friend request exists failed: %w", err)
	}
	return count > 0, nil
}

type GroupJoinRequestRepo struct {
	db *gorm.DB
}

func NewGroupJoinRequestRepo(db *gorm.DB) *GroupJoinRequestRepo {
	return &GroupJoinRequestRepo{db: db}
}

func (r *GroupJoinRequestRepo) Create(ctx context.Context, req *model.GroupJoinRequest) error {
	return r.db.WithContext(ctx).Create(req).Error
}

func (r *GroupJoinRequestRepo) GetByID(ctx context.Context, reqID string) (*model.GroupJoinRequest, error) {
	var req model.GroupJoinRequest
	err := r.db.WithContext(ctx).Where("id = ?", reqID).First(&req).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrRequestNotFound
		}
		return nil, fmt.Errorf("query group join request failed: %w", err)
	}
	return &req, nil
}

func (r *GroupJoinRequestRepo) GetPendingRequests(ctx context.Context, gid string) ([]*model.GroupJoinRequest, error) {
	var requests []*model.GroupJoinRequest
	err := r.db.WithContext(ctx).
		Where("group_id = ? AND status = ?", gid, model.GroupJoinRequestPending).
		Order("created_at DESC").
		Find(&requests).Error
	if err != nil {
		return nil, fmt.Errorf("query pending group join requests failed: %w", err)
	}
	return requests, nil
}

func (r *GroupJoinRequestRepo) UpdateStatus(ctx context.Context, reqID, status string) error {
	result := r.db.WithContext(ctx).
		Model(&model.GroupJoinRequest{}).
		Where("id = ?", reqID).
		Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("update group join request status failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrRequestNotFound
	}
	return nil
}

func (r *GroupJoinRequestRepo) Delete(ctx context.Context, reqID string) error {
	result := r.db.WithContext(ctx).Delete(&model.GroupJoinRequest{}, "id = ?", reqID)
	if result.Error != nil {
		return fmt.Errorf("delete group join request failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrRequestNotFound
	}
	return nil
}

func (r *GroupJoinRequestRepo) Exists(ctx context.Context, userID, groupID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.GroupJoinRequest{}).
		Where("user_id = ? AND group_id = ? AND status = ?", userID, groupID, model.GroupJoinRequestPending).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check group join request exists failed: %w", err)
	}
	return count > 0, nil
}
