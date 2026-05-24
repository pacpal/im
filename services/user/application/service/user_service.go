// Package service 提供 user 服务的业务逻辑实现（Use case 层）。
package service

import (
	"IM/pkg/auth"
	"IM/pkg/id"
	"IM/pkg/logger"
	"IM/services/user/domain/entity"
	"IM/services/user/domain/event"
	"IM/services/user/domain/repository"
	"IM/services/user/infrastructure/cache"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	// 常见的业务错误，供上层调用者判断处理。
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrAlreadyFriends    = errors.New("already friends")
	ErrNotFriends        = errors.New("not friends")
	ErrRequestNotFound   = errors.New("request not found")
	ErrRequestExists     = errors.New("request already exists")
	ErrCannotAddSelf     = errors.New("cannot add yourself as friend")
	ErrInvalidRequest    = errors.New("invalid request")
)

// UserService 实现了用户相关的业务逻辑（注册、登录、加好友、处理请求等）。
type UserService struct {
	userRepo          repository.UserRepository
	friendshipRepo    repository.FriendshipRepository
	friendRequestRepo repository.FriendRequestRepository
	userCache         *cache.UserCache
	idGenerator       *id.SnowflakeGenerator
	jwtUtil           *auth.JWTUtil
	eventPublisher    *event.EventPublisher
}

// NewUserService 构造 UserService。
func NewUserService(
	userRepo repository.UserRepository,
	friendshipRepo repository.FriendshipRepository,
	friendRequestRepo repository.FriendRequestRepository,
	userCache *cache.UserCache,
	idGenerator *id.SnowflakeGenerator,
	jwtUtil *auth.JWTUtil,
	eventPublisher *event.EventPublisher,
) *UserService {
	return &UserService{
		userRepo:          userRepo,
		friendshipRepo:    friendshipRepo,
		friendRequestRepo: friendRequestRepo,
		userCache:         userCache,
		idGenerator:       idGenerator,
		jwtUtil:           jwtUtil,
		eventPublisher:    eventPublisher,
	}
}

func (s *UserService) Register(ctx context.Context, tele, name, password string) (res *entity.User, err error) {
	done := logger.StartStep("UserService.Register", "tele", tele)
	defer func() { done(err) }()

	var exists bool
	exists, err = s.userRepo.ExistsByTele(ctx, tele)
	if err != nil {
		return nil, err
	}
	if exists {
		err = ErrUserAlreadyExists
		return nil, err
	}

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userID := s.idGenerator.Generate()
	res = entity.NewUser(userID, name, tele, string(hashedPassword))

	if err = s.userRepo.Create(ctx, res); err != nil {
		return nil, err
	}

	s.eventPublisher.Publish(&event.UserRegisteredEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "user.registered",
			OccurredAt:  time.Now(),
			AggregateID: res.ID,
		},
		UserID: res.ID,
		Name:   res.Name,
		Tele:   res.Tele,
	})

	logger.Infow("Register: user created", "component", "user_service", "user_id", res.ID)
	return res, nil
}

func (s *UserService) Login(ctx context.Context, tele, id, password string) (res *entity.User, token string, err error) {
	done := logger.StartStep("UserService.Login", "tele", tele, "id", id)
	defer func() { done(err) }()

	if tele != "" {
		res, err = s.userRepo.GetByTele(ctx, tele)
	} else if id != "" {
		res, err = s.userRepo.GetByID(ctx, id)
	} else {
		err = ErrUserNotFound
		return
	}
	if err != nil {
		err = ErrUserNotFound
		return
	}

	if e := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password)); e != nil {
		err = ErrInvalidPassword
		return
	}

	token, err = s.jwtUtil.GenerateToken(res.ID, res.Name)
	if err != nil {
		return
	}

	s.eventPublisher.Publish(&event.UserLoggedInEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "user.logged_in",
			OccurredAt:  time.Now(),
			AggregateID: res.ID,
		},
		UserID: res.ID,
		Tele:   res.Tele,
	})

	logger.Infow("Login: user logged in", "component", "user_service", "user_id", res.ID)
	return
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (res *entity.User, err error) {
	done := logger.StartStep("UserService.GetUserByID", "id", id)
	defer func() { done(err) }()

	res, err = s.userRepo.GetByID(ctx, id)
	if err != nil {
		err = ErrUserNotFound
		return
	}
	logger.Infow("GetUserByID: found", "component", "user_service", "user_id", res.ID)
	return
}

func (s *UserService) GetUserByTele(ctx context.Context, tele string) (res *entity.User, err error) {
	done := logger.StartStep("UserService.GetUserByTele", "tele", tele)
	defer func() { done(err) }()

	res, err = s.userRepo.GetByTele(ctx, tele)
	if err != nil {
		err = ErrUserNotFound
		return
	}
	logger.Infow("GetUserByTele: found", "component", "user_service", "user_id", res.ID)
	return
}

func (s *UserService) GetUsersByIDs(ctx context.Context, ids []string) (res []*entity.User, err error) {
	done := logger.StartStep("UserService.GetUsersByIDs", "count", len(ids))
	defer func() { done(err) }()

	res, err = s.userRepo.GetByIDs(ctx, ids)
	return
}

func (s *UserService) UpdateUser(ctx context.Context, user *entity.User) (err error) {
	done := logger.StartStep("UserService.UpdateUser", "user_id", user.ID)
	defer func() { done(err) }()

	user.UpdatedAt = time.Now()
	err = s.userRepo.Update(ctx, user)
	if err == nil {
		logger.Infow("UpdateUser: updated", "component", "user_service", "user_id", user.ID)
	}
	return
}

func (s *UserService) GetFriends(ctx context.Context, userID string) (res []*entity.User, err error) {
	done := logger.StartStep("UserService.GetFriends", "user_id", userID)
	defer func() { done(err) }()

	res, err = s.friendshipRepo.GetFriends(ctx, userID)
	if err == nil {
		logger.Infow("GetFriends: retrieved", "component", "user_service", "user_id", userID, "count", len(res))
	}
	return
}

func (s *UserService) AddFriend(ctx context.Context, fromUID, toUID, reason string) error {
	done := logger.StartStep("UserService.AddFriend", "from", fromUID, "to", toUID)

	if fromUID == toUID {
		done(ErrCannotAddSelf)
		return ErrCannotAddSelf
	}

	_, err := s.userRepo.GetByID(ctx, toUID)
	if err != nil {
		done(err)
		return ErrUserNotFound
	}
	logger.Infow("AddFriend: target user found", "component", "user_service", "from", fromUID, "to", toUID)

	exists, err := s.friendshipRepo.Exists(ctx, fromUID, toUID)
	if err != nil {
		done(err)
		return err
	}
	if exists {
		done(ErrAlreadyFriends)
		logger.Infow("AddFriend: already friends", "component", "user_service", "from", fromUID, "to", toUID)
		return ErrAlreadyFriends
	}

	requestExists, err := s.friendRequestRepo.Exists(ctx, fromUID, toUID)
	if err != nil {
		done(err)
		return err
	}
	if requestExists {
		done(ErrRequestExists)
		logger.Infow("AddFriend: request already exists", "component", "user_service", "from", fromUID, "to", toUID)
		return ErrRequestExists
	}

	requestID := s.idGenerator.Generate()
	friendRequest := entity.NewFriendRequest(requestID, fromUID, toUID, reason)

	if err := s.friendRequestRepo.Create(ctx, friendRequest); err != nil {
		done(err)
		return err
	}
	logger.Infow("AddFriend: friend request created", "component", "user_service", "request_id", requestID, "from", fromUID, "to", toUID)

	s.eventPublisher.Publish(&event.FriendRequestCreatedEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "friend_request.created",
			OccurredAt:  time.Now(),
			AggregateID: requestID,
		},
		RequestID: requestID,
		FromUID:   fromUID,
		ToUID:     toUID,
	})

	logger.Infow("AddFriend: published FriendRequestCreatedEvent", "component", "user_service", "request_id", requestID)
	done(nil)
	return nil
}

func (s *UserService) AcceptFriendRequest(ctx context.Context, requestID, acceptorID string) (err error) {
	done := logger.StartStep("UserService.AcceptFriendRequest", "request_id", requestID, "acceptor", acceptorID)
	defer func() { done(err) }()

	var req *entity.FriendRequest
	req, err = s.friendRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		err = ErrRequestNotFound
		return
	}

	if req.ToUID != acceptorID {
		err = ErrInvalidRequest
		return
	}

	if !req.IsPending() {
		err = ErrInvalidRequest
		return
	}

	friendship1 := entity.NewFriendship(req.FromUID, req.ToUID)
	friendship2 := entity.NewFriendship(req.ToUID, req.FromUID)

	if err = s.friendshipRepo.Create(ctx, friendship1); err != nil {
		return
	}
	if err = s.friendshipRepo.Create(ctx, friendship2); err != nil {
		return
	}

	req.Accept()
	if err = s.friendRequestRepo.UpdateStatus(ctx, requestID, entity.FriendRequestStatusAccepted); err != nil {
		return
	}

	s.eventPublisher.Publish(&event.FriendRequestAcceptedEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "friend_request.accepted",
			OccurredAt:  time.Now(),
			AggregateID: requestID,
		},
		RequestID: requestID,
		FromUID:   req.FromUID,
		ToUID:     req.ToUID,
	})

	logger.Infow("AcceptFriendRequest: accepted", "component", "user_service", "request_id", requestID, "from", req.FromUID, "to", req.ToUID)
	return
}

func (s *UserService) RejectFriendRequest(ctx context.Context, requestID, rejecterID string) (err error) {
	done := logger.StartStep("UserService.RejectFriendRequest", "request_id", requestID, "rejecter", rejecterID)
	defer func() { done(err) }()

	var req *entity.FriendRequest
	req, err = s.friendRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		err = ErrRequestNotFound
		return
	}

	if req.ToUID != rejecterID {
		err = ErrInvalidRequest
		return
	}

	if !req.IsPending() {
		err = ErrInvalidRequest
		return
	}

	req.Reject()
	if err = s.friendRequestRepo.UpdateStatus(ctx, requestID, entity.FriendRequestStatusRejected); err != nil {
		return
	}

	s.eventPublisher.Publish(&event.FriendRequestRejectedEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "friend_request.rejected",
			OccurredAt:  time.Now(),
			AggregateID: requestID,
		},
		RequestID: requestID,
		FromUID:   req.FromUID,
		ToUID:     req.ToUID,
	})

	logger.Infow("RejectFriendRequest: rejected", "component", "user_service", "request_id", requestID, "from", req.FromUID, "to", req.ToUID)
	return
}

func (s *UserService) GetPendingFriendRequests(ctx context.Context, uid string) (res []*entity.FriendRequest, err error) {
	done := logger.StartStep("UserService.GetPendingFriendRequests", "user_id", uid)
	defer func() { done(err) }()

	res, err = s.friendRequestRepo.GetPendingRequests(ctx, uid)
	if err == nil {
		logger.Infow("GetPendingFriendRequests: retrieved", "component", "user_service", "user_id", uid, "count", len(res))
	}
	return
}

func (s *UserService) RemoveFriend(ctx context.Context, userID, friendID string) (err error) {
	done := logger.StartStep("UserService.RemoveFriend", "user_id", userID, "friend_id", friendID)
	defer func() { done(err) }()

	if err = s.friendshipRepo.Delete(ctx, userID, friendID); err != nil {
		return
	}
	if err = s.friendshipRepo.Delete(ctx, friendID, userID); err != nil {
		return
	}

	s.eventPublisher.Publish(&event.FriendshipDeletedEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "friendship.deleted",
			OccurredAt:  time.Now(),
			AggregateID: userID,
		},
		UserID:   userID,
		FriendID: friendID,
	})

	logger.Infow("RemoveFriend: removed", "component", "user_service", "user_id", userID, "friend_id", friendID)
	return
}

func (s *UserService) CheckFriendship(ctx context.Context, userID1, userID2 string) (ok bool, err error) {
	done := logger.StartStep("UserService.CheckFriendship", "user1", userID1, "user2", userID2)
	defer func() { done(err) }()

	ok, err = s.friendshipRepo.Exists(ctx, userID1, userID2)
	if err == nil {
		logger.Infow("CheckFriendship: result", "component", "user_service", "user1", userID1, "user2", userID2, "ok", ok)
	}
	return
}
