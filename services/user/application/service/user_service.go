package service

import (
	"IM/pkg/auth"
	"IM/pkg/id"
	"IM/services/user/domain/entity"
	"IM/services/user/domain/event"
	"IM/services/user/domain/repository"
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
	ErrInvalidRequest     = errors.New("invalid request")
)

type UserService struct {
	userRepo         repository.UserRepository
	friendshipRepo   repository.FriendshipRepository
	friendRequestRepo repository.FriendRequestRepository
	userCache        *UserCache
	idGenerator      *id.SnowflakeGenerator
	jwtUtil          *auth.JWTUtil
	eventPublisher   *event.EventPublisher
}

type UserCache struct {
	getUserFunc func(ctx context.Context, userID string) (*entity.User, error)
}

func NewUserCache(getUserFunc func(ctx context.Context, userID string) (*entity.User, error)) *UserCache {
	return &UserCache{getUserFunc: getUserFunc}
}

func NewUserService(
	userRepo repository.UserRepository,
	friendshipRepo repository.FriendshipRepository,
	friendRequestRepo repository.FriendRequestRepository,
	userCache *UserCache,
	idGenerator *id.SnowflakeGenerator,
	jwtUtil *auth.JWTUtil,
	eventPublisher *event.EventPublisher,
) *UserService {
	return &UserService{
		userRepo:         userRepo,
		friendshipRepo:   friendshipRepo,
		friendRequestRepo: friendRequestRepo,
		userCache:        userCache,
		idGenerator:      idGenerator,
		jwtUtil:          jwtUtil,
		eventPublisher:   eventPublisher,
	}
}

func (s *UserService) Register(ctx context.Context, tele, name, password string) (*entity.User, error) {
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

	userID := s.idGenerator.Generate()
	user := entity.NewUser(userID, name, tele, string(hashedPassword))

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	s.eventPublisher.Publish(&event.UserRegisteredEvent{
		BaseEvent: event.BaseEvent{
			eventType:   "user.registered",
			occurredAt:  time.Now(),
			aggregateID: user.ID,
		},
		UserID: user.ID,
		Name:   user.Name,
		Tele:   user.Tele,
	})

	return user, nil
}

func (s *UserService) Login(ctx context.Context, tele, password string) (*entity.User, string, error) {
	user, err := s.userRepo.GetByTele(ctx, tele)
	if err != nil {
		return nil, "", ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", ErrInvalidPassword
	}

	token, err := s.jwtUtil.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	s.eventPublisher.Publish(&event.UserLoggedInEvent{
		BaseEvent: event.BaseEvent{
			eventType:   "user.logged_in",
			occurredAt:  time.Now(),
			aggregateID: user.ID,
		},
		UserID: user.ID,
		Tele:   user.Tele,
	})

	return user, token, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) GetUserByTele(ctx context.Context, tele string) (*entity.User, error) {
	user, err := s.userRepo.GetByTele(ctx, tele)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) GetUsersByIDs(ctx context.Context, ids []string) ([]*entity.User, error) {
	return s.userRepo.GetByIDs(ctx, ids)
}

func (s *UserService) UpdateUser(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) GetFriends(ctx context.Context, userID string) ([]*entity.User, error) {
	return s.friendshipRepo.GetFriends(ctx, userID)
}

func (s *UserService) AddFriend(ctx context.Context, fromUID, toUID, reason string) error {
	if fromUID == toUID {
		return ErrCannotAddSelf
	}

	_, err := s.userRepo.GetByID(ctx, toUID)
	if err != nil {
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

	requestID := s.idGenerator.Generate()
	friendRequest := entity.NewFriendRequest(requestID, fromUID, toUID, reason)

	if err := s.friendRequestRepo.Create(ctx, friendRequest); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.FriendRequestCreatedEvent{
		BaseEvent: event.BaseEvent{
			eventType:   "friend_request.created",
			occurredAt:  time.Now(),
			aggregateID: requestID,
		},
		RequestID: requestID,
		FromUID:   fromUID,
		ToUID:     toUID,
	})

	return nil
}

func (s *UserService) AcceptFriendRequest(ctx context.Context, requestID, acceptorID string) error {
	req, err := s.friendRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		return ErrRequestNotFound
	}

	if req.ToUID != acceptorID {
		return ErrInvalidRequest
	}

	if !req.IsPending() {
		return ErrInvalidRequest
	}

	friendship1 := entity.NewFriendship(req.FromUID, req.ToUID)
	friendship2 := entity.NewFriendship(req.ToUID, req.FromUID)

	if err := s.friendshipRepo.Create(ctx, friendship1); err != nil {
		return err
	}
	if err := s.friendshipRepo.Create(ctx, friendship2); err != nil {
		return err
	}

	req.Accept()
	if err := s.friendRequestRepo.UpdateStatus(ctx, requestID, entity.FriendRequestStatusAccepted); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.FriendRequestAcceptedEvent{
		BaseEvent: event.BaseEvent{
			eventType:   "friend_request.accepted",
			occurredAt:  time.Now(),
			aggregateID: requestID,
		},
		RequestID: requestID,
		FromUID:   req.FromUID,
		ToUID:     req.ToUID,
	})

	return nil
}

func (s *UserService) RejectFriendRequest(ctx context.Context, requestID, rejecterID string) error {
	req, err := s.friendRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		return ErrRequestNotFound
	}

	if req.ToUID != rejecterID {
		return ErrInvalidRequest
	}

	if !req.IsPending() {
		return ErrInvalidRequest
	}

	req.Reject()
	if err := s.friendRequestRepo.UpdateStatus(ctx, requestID, entity.FriendRequestStatusRejected); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.FriendRequestRejectedEvent{
		BaseEvent: event.BaseEvent{
			eventType:   "friend_request.rejected",
			occurredAt:  time.Now(),
			aggregateID: requestID,
		},
		RequestID: requestID,
		FromUID:   req.FromUID,
		ToUID:     req.ToUID,
	})

	return nil
}

func (s *UserService) GetPendingFriendRequests(ctx context.Context, uid string) ([]*entity.FriendRequest, error) {
	return s.friendRequestRepo.GetPendingRequests(ctx, uid)
}

func (s *UserService) RemoveFriend(ctx context.Context, userID, friendID string) error {
	if err := s.friendshipRepo.Delete(ctx, userID, friendID); err != nil {
		return err
	}
	if err := s.friendshipRepo.Delete(ctx, friendID, userID); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.FriendshipDeletedEvent{
		BaseEvent: event.BaseEvent{
			eventType:   "friendship.deleted",
			occurredAt:  time.Now(),
			aggregateID: userID,
		},
		UserID:   userID,
		FriendID: friendID,
	})

	return nil
}

func (s *UserService) CheckFriendship(ctx context.Context, userID1, userID2 string) (bool, error) {
	return s.friendshipRepo.Exists(ctx, userID1, userID2)
}
