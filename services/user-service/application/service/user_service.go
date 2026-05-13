package service

import (
	"IM/services/user-service/domain/entity"
	"IM/services/user-service/infrastructure/repository"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrAlreadyFriends     = errors.New("already friends")
	ErrNotFriends         = errors.New("not friends")
	ErrRequestNotFound    = errors.New("request not found")
	ErrRequestExists      = errors.New("request already exists")
	ErrCannotAddSelf      = errors.New("cannot add yourself as friend")
)

type UserApplicationService struct {
	userRepo         repository.UserRepository
	friendshipRepo   repository.FriendshipRepository
	friendRequestRepo repository.FriendRequestRepository
}

func NewUserApplicationService(
	userRepo repository.UserRepository,
	friendshipRepo repository.FriendshipRepository,
	friendRequestRepo repository.FriendRequestRepository,
) *UserApplicationService {
	return &UserApplicationService{
		userRepo:         userRepo,
		friendshipRepo:   friendshipRepo,
		friendRequestRepo: friendRequestRepo,
	}
}

func (s *UserApplicationService) Register(ctx context.Context, tele, name, password string) (*entity.User, error) {
	exists, err := s.userRepo.ExistsByTele(ctx, tele)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:        generateID("user"),
		Name:      name,
		Tele:      tele,
		Password:  string(hashedPassword),
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserApplicationService) Login(ctx context.Context, tele, password string) (*entity.User, error) {
	user, err := s.userRepo.GetByTele(ctx, tele)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

func (s *UserApplicationService) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserApplicationService) GetUserByTele(ctx context.Context, tele string) (*entity.User, error) {
	user, err := s.userRepo.GetByTele(ctx, tele)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserApplicationService) GetFriends(ctx context.Context, userID string) ([]*entity.User, error) {
	return s.friendshipRepo.GetFriends(ctx, userID)
}

func (s *UserApplicationService) AddFriend(ctx context.Context, fromUID, toUID, reason string) error {
	if fromUID == toUID {
		return ErrCannotAddSelf
	}

	toUser, err := s.userRepo.GetByID(ctx, toUID)
	if err != nil {
		return ErrUserNotFound
	}
	if toUser == nil {
		return ErrUserNotFound
	}

	exists, err := s.friendshipRepo.Exists(ctx, fromUID, toUID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyFriends
	}

	requestExists, err := s.friendRequestRepo.Exists(ctx, fromUID, toUID)
	if err != nil {
		return err
	}
	if requestExists {
		return ErrRequestExists
	}

	friendRequest := &entity.FriendRequest{
		ID:        generateID("freq"),
		FromUID:   fromUID,
		ToUID:     toUID,
		Reason:    reason,
		Status:    entity.FriendRequestStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.friendRequestRepo.Create(ctx, friendRequest)
}

func (s *UserApplicationService) AcceptFriendRequest(ctx context.Context, requestID, acceptorID string) error {
	req, err := s.friendRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		return ErrRequestNotFound
	}
	if req.ToUID != acceptorID {
		return ErrRequestNotFound
	}

	friendship1 := &entity.Friendship{
		UserID:    req.FromUID,
		FriendID: req.ToUID,
		Status:    entity.FriendshipStatusActive,
		CreatedAt: time.Now(),
	}

	friendship2 := &entity.Friendship{
		UserID:    req.ToUID,
		FriendID: req.FromUID,
		Status:    entity.FriendshipStatusActive,
		CreatedAt: time.Now(),
	}

	if err := s.friendshipRepo.Create(ctx, friendship1); err != nil {
		return err
	}
	if err := s.friendshipRepo.Create(ctx, friendship2); err != nil {
		return err
	}

	return s.friendRequestRepo.UpdateStatus(ctx, requestID, entity.FriendRequestStatusAccepted)
}

func (s *UserApplicationService) RejectFriendRequest(ctx context.Context, requestID, rejecterID string) error {
	req, err := s.friendRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		return ErrRequestNotFound
	}
	if req.ToUID != rejecterID {
		return ErrRequestNotFound
	}

	return s.friendRequestRepo.UpdateStatus(ctx, requestID, entity.FriendRequestStatusRejected)
}

func (s *UserApplicationService) GetPendingFriendRequests(ctx context.Context, uid string) ([]*entity.FriendRequest, error) {
	return s.friendRequestRepo.GetPendingRequests(ctx, uid)
}

func (s *UserApplicationService) RemoveFriend(ctx context.Context, userID, friendID string) error {
	if err := s.friendshipRepo.Delete(ctx, userID, friendID); err != nil {
		return err
	}
	return s.friendshipRepo.Delete(ctx, friendID, userID)
}

func (s *UserApplicationService) UpdateUser(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(ctx, user)
}

func generateID(prefix string) string {
	return prefix + "_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}