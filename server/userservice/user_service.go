package userservice

import (
	"IM/server/model"
	repo "IM/server/repository"
	"context"

	"golang.org/x/crypto/bcrypt"
)

const jwtSecret = "jwtsecret"

type UserService struct {
	userRepo  repo.UserRepo
	groupRepo repo.GroupRepo
}

func NewUserService(ur repo.UserRepo, gr repo.GroupRepo) *UserService {
	return &UserService{userRepo: ur, groupRepo: gr}
}
func (s *UserService) Register(ctx context.Context, tele, password string) (*model.User, error) {
	_, err := s.userRepo.GetUserByTele(ctx, tele)
	if err == nil {
		return nil, model.ErrUserExists
	}
	secretWord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	uid := "newuid"
	user := &model.User{
		ID:       uid,
		Password: string(secretWord),
		Tele:     tele,
	}
	if err := s.userRepo.RefreshUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserService) Login(ctx context.Context, id, tele, password string) error {
	if tele == "" {
		user, err := s.userRepo.GetUserByID(ctx, id)
		if err != nil {
			return err
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return err
		}
	} else if id == "" {
		user, err := s.userRepo.GetUserByTele(ctx, tele)
		if err != nil {
			return err
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return err
		}
	}
	return nil
}
func (s *UserService) GetFriends(ctx context.Context, id string) ([]string, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	//可以选择分页限制
	friends := make([]string, 0, len(user.Friends))
	for id, ok := range user.Friends {
		if ok {
			friends = append(friends, id)
		}
	}
	return friends, nil
}
func (s *UserService) AddFriend(ctx context.Context, uid, targetID, reason string) error {
	_, err := s.userRepo.GetUserByID(ctx, uid)
	if err != nil {
		return err
	}
	_, err = s.userRepo.GetUserByID(ctx, targetID)
	if err != nil {
		return err
	}
	//logic add later...notify
	return nil
}

func (s *UserService) RemoveFriend(ctx context.Context, uid, targetID string) error {
	user, err := s.userRepo.GetUserByID(ctx, uid)
	if err != nil {
		return err
	}
	if err := user.RemoveFriend(targetID); err != nil {
		return err
	}
	return s.userRepo.RefreshUser(ctx, user)
}

func (s *UserService) ReplyFriendAdd(ctx context.Context, uid, targetID, reply string) error {
	if reply != "agree" {
		return nil
	}
	user, err := s.userRepo.GetUserByID(ctx, uid)
	if err != nil {
		return err
	}
	if err := user.AddFriend(targetID); err != nil {
		return err
	}
	if err := s.userRepo.RefreshUser(ctx, user); err != nil {
		return err
	}

	target, err := s.userRepo.GetUserByID(ctx, targetID)
	if err != nil {
		return err
	}
	if err := target.AddFriend(uid); err != nil {
		return err
	}
	return s.userRepo.RefreshUser(ctx, target)
}

func (s *UserService) CheckFriendship(ctx context.Context, uid1, uid2 string) (bool, error) {
	user, err := s.userRepo.GetUserByID(ctx, uid1)
	if err != nil {
		return false, err
	}
	return user.Friends[uid2], nil
}
