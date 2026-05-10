// Package userservice for normal user oprate and relationship...后续优化
package userservice

import (
	"IM/server/model"
	repo "IM/server/repository"
	"context"
)

type GroupService struct {
	userRepo  repo.UserRepo
	groupRepo repo.GroupRepo
}

func NewGroupAppService(ur repo.UserRepo, gr repo.GroupRepo) *GroupService {
	return &GroupService{
		userRepo:  ur,
		groupRepo: gr,
	}
}

func (s *GroupService) CreateGroup(ctx context.Context, ownerID, name string) (*model.Group, error) {
	_, err := s.userRepo.GetUserByID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	gid := "GenGroupID()"
	group := model.NewGroup(gid, name, ownerID)
	if err := s.groupRepo.SaveGroup(ctx, group); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByID(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	user.JoinGroup(gid)
	if err := s.userRepo.RefreshUser(ctx, user); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *GroupService) GetGroup(ctx context.Context, gid string) (*model.Group, error) {
	return s.groupRepo.GetGroupByID(ctx, gid)
}

func (s *GroupService) JoinGroup(ctx context.Context, uid, gid, reason string) error {
	_, err := s.userRepo.GetUserByID(ctx, uid)
	if err != nil {
		return err
	}
	_, err = s.groupRepo.GetGroupByID(ctx, gid)
	if err != nil {
		return err
	}
	return nil
}

func (s *GroupService) LeaveGroup(ctx context.Context, uid, gid string) error {
	user, err := s.userRepo.GetUserByID(ctx, uid)
	if err != nil {
		return err
	}
	if err := user.LeaveGroup(gid); err != nil {
		return err
	}
	if err := s.userRepo.RefreshUser(ctx, user); err != nil {
		return err
	}

	group, err := s.groupRepo.GetGroupByID(ctx, gid)
	if err != nil {
		return err
	}
	if err := group.RemoveMember(uid); err != nil {
		return err
	}
	return s.groupRepo.SaveGroup(ctx, group)
}

func (s *GroupService) ReplyGroupAdd(ctx context.Context, ownerID, uid, gid, reply string) error {
	if reply != "agree" {
		return nil
	}
	group, err := s.groupRepo.GetGroupByID(ctx, gid)
	if err != nil {
		return err
	}
	if group.OwnerID != ownerID {
		return model.ErrNotMember
	}
	if err := group.AddMember(uid); err != nil {
		return err
	}
	if err := s.groupRepo.SaveGroup(ctx, group); err != nil {
		return err
	}

	user, err := s.userRepo.GetUserByID(ctx, uid)
	if err != nil {
		return err
	}
	user.JoinGroup(gid)
	return s.userRepo.RefreshUser(ctx, user)
}

func (s *GroupService) IsGroupMember(ctx context.Context, gid, uid string) (bool, error) {
	group, err := s.groupRepo.GetGroupByID(ctx, gid)
	if err != nil {
		return false, err
	}
	return group.IsMember(uid), nil
}
