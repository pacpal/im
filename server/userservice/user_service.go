// Package userservice 提供用户相关的业务逻辑服务，包括用户注册、登录、好友管理等功能。
package userservice

import (
	"IM/server/model"
	repo "IM/server/repository"
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo          repo.UserRepo
	friendshipRepo    repo.FriendshipRepo
	friendRequestRepo repo.FriendRequestRepo
}

func NewUserService(ur repo.UserRepo, fr repo.FriendshipRepo, frr repo.FriendRequestRepo) *UserService {
	return &UserService{
		userRepo:          ur,
		friendshipRepo:    fr,
		friendRequestRepo: frr,
	}
}

func (s *UserService) Register(ctx context.Context, tele, name, password string) (*model.User, error) {
	exists, err := s.userRepo.ExistsByTele(ctx, tele)
	if err != nil {
		return nil, fmt.Errorf("check user existence failed: %w", err)
	}
	if exists {
		return nil, model.ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password failed: %w", err)
	}

	uid := fmt.Sprintf("user_%d", time.Now().UnixNano())
	user := &model.User{
		ID:       uid,
		Name:     name,
		Tele:     tele,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user failed: %w", err)
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, tele, password string) (*model.User, error) {
	user, err := s.userRepo.GetByTele(ctx, tele)
	if err != nil {
		return nil, fmt.Errorf("get user failed: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, model.ErrInvalidPassword
	}

	return user, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, uid string) (*model.User, error) {
	return s.userRepo.GetByID(ctx, uid)
}

func (s *UserService) GetFriends(ctx context.Context, uid string) ([]*model.User, error) {
	return s.friendshipRepo.GetFriends(ctx, uid)
}

func (s *UserService) AddFriend(ctx context.Context, uid, targetID, reason string) error {
	if uid == targetID {
		return fmt.Errorf("cannot add yourself as friend")
	}

	exists, err := s.userRepo.Exists(ctx, targetID)
	if err != nil {
		return fmt.Errorf("check target user failed: %w", err)
	}
	if !exists {
		return model.ErrUserNotFound
	}

	exists, err = s.friendshipRepo.Exists(ctx, uid, targetID)
	if err != nil {
		return fmt.Errorf("check friendship failed: %w", err)
	}
	if exists {
		return model.ErrAlreadyFriend
	}

	exists, err = s.friendRequestRepo.Exists(ctx, uid, targetID)
	if err != nil {
		return fmt.Errorf("check friend request failed: %w", err)
	}
	if exists {
		return fmt.Errorf("friend request already sent")
	}

	requestID := fmt.Sprintf("freq_%d", time.Now().UnixNano())
	friendReq := &model.FriendRequest{
		ID:      requestID,
		FromUID: uid,
		ToUID:   targetID,
		Reason:  reason,
		Status:  model.FriendRequestPending,
	}

	return s.friendRequestRepo.Create(ctx, friendReq)
}

func (s *UserService) RemoveFriend(ctx context.Context, uid, targetID string) error {
	if uid == targetID {
		return fmt.Errorf("cannot remove yourself")
	}

	exists, err := s.friendshipRepo.Exists(ctx, uid, targetID)
	if err != nil {
		return fmt.Errorf("check friendship failed: %w", err)
	}
	if !exists {
		return model.ErrNotFriend
	}

	if err := s.friendshipRepo.Delete(ctx, uid, targetID); err != nil {
		return fmt.Errorf("remove friend failed: %w", err)
	}

	return nil
}

func (s *UserService) ReplyFriendAdd(ctx context.Context, uid, requestID, reply string) error {
	request, err := s.friendRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		return err
	}

	if request.ToUID != uid {
		return model.ErrInvalidRequest
	}

	if request.Status != model.FriendRequestPending {
		return fmt.Errorf("request already processed")
	}

	if reply != "agree" {
		return s.friendRequestRepo.UpdateStatus(ctx, requestID, model.FriendRequestRejected)
	}

	friendship1 := &model.Friendship{
		UserID:   uid,
		FriendID: request.FromUID,
		Status:   1,
	}

	friendship2 := &model.Friendship{
		UserID:   request.FromUID,
		FriendID: uid,
		Status:   1,
	}

	if err := s.friendshipRepo.Create(ctx, friendship1); err != nil {
		return fmt.Errorf("create friendship failed: %w", err)
	}

	if err := s.friendshipRepo.Create(ctx, friendship2); err != nil {
		_ = s.friendshipRepo.Delete(ctx, uid, request.FromUID)
		return fmt.Errorf("create friendship failed: %w", err)
	}

	return s.friendRequestRepo.UpdateStatus(ctx, requestID, model.FriendRequestAccepted)
}

func (s *UserService) GetPendingFriendRequests(ctx context.Context, uid string) ([]*model.FriendRequest, error) {
	return s.friendRequestRepo.GetPendingRequests(ctx, uid)
}

func (s *UserService) CheckFriendship(ctx context.Context, uid1, uid2 string) (bool, error) {
	return s.friendshipRepo.Exists(ctx, uid1, uid2)
}
