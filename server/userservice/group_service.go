package userservice

import (
	"IM/server/model"
	repo "IM/server/repository"
	"context"
	"fmt"
	"time"
)

type GroupService struct {
	groupRepo        repo.GroupRepo
	groupMemberRepo  repo.GroupMemberRepo
	userRepo         repo.UserRepo
	groupJoinRequest repo.GroupJoinRequestRepo
}

func NewGroupService(gr repo.GroupRepo, gmr repo.GroupMemberRepo, ur repo.UserRepo, gjr repo.GroupJoinRequestRepo) *GroupService {
	return &GroupService{
		groupRepo:        gr,
		groupMemberRepo:  gmr,
		userRepo:         ur,
		groupJoinRequest: gjr,
	}
}

func (s *GroupService) CreateGroup(ctx context.Context, ownerID, name, description string) (*model.Group, error) {
	_, err := s.userRepo.GetByID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	gid := fmt.Sprintf("group_%d", time.Now().UnixNano())
	group := &model.Group{
		ID:          gid,
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		Type:        "normal",
	}

	if err := s.groupRepo.Create(ctx, group); err != nil {
		return nil, err
	}

	groupMember := &model.GroupMember{
		GroupID: gid,
		UserID:  ownerID,
		Role:    3,
	}

	if err := s.groupMemberRepo.AddMember(ctx, groupMember); err != nil {
		_ = s.groupRepo.Delete(ctx, gid)
		return nil, fmt.Errorf("add owner to group failed: %w", err)
	}

	return group, nil
}

func (s *GroupService) GetGroup(ctx context.Context, gid string) (*model.Group, error) {
	return s.groupRepo.GetByID(ctx, gid)
}

func (s *GroupService) JoinGroup(ctx context.Context, uid, gid, reason string) error {
	_, err := s.groupRepo.GetByID(ctx, gid)
	if err != nil {
		return err
	}

	isMember, err := s.groupMemberRepo.IsMember(ctx, gid, uid)
	if err != nil {
		return err
	}
	if isMember {
		return model.ErrAlreadyMember
	}

	exists, err := s.groupJoinRequest.Exists(ctx, uid, gid)
	if err != nil {
		return fmt.Errorf("check join request failed: %w", err)
	}
	if exists {
		return fmt.Errorf("join request already sent")
	}

	requestID := fmt.Sprintf("greq_%d", time.Now().UnixNano())
	joinReq := &model.GroupJoinRequest{
		ID:      requestID,
		UserID:  uid,
		GroupID: gid,
		Reason:  reason,
		Status:  model.GroupJoinRequestPending,
	}

	return s.groupJoinRequest.Create(ctx, joinReq)
}

func (s *GroupService) LeaveGroup(ctx context.Context, uid, gid string) error {
	group, err := s.groupRepo.GetByID(ctx, gid)
	if err != nil {
		return err
	}

	if group.OwnerID == uid {
		return fmt.Errorf("owner cannot leave the group")
	}

	isMember, err := s.groupMemberRepo.IsMember(ctx, gid, uid)
	if err != nil {
		return err
	}
	if !isMember {
		return model.ErrNotMember
	}

	return s.groupMemberRepo.RemoveMember(ctx, gid, uid)
}

func (s *GroupService) ReplyGroupAdd(ctx context.Context, ownerID, gid, reply string) error {
	group, err := s.groupRepo.GetByID(ctx, gid)
	if err != nil {
		return err
	}

	if group.OwnerID != ownerID {
		return model.ErrNotOwner
	}

	requests, err := s.groupJoinRequest.GetPendingRequests(ctx, gid)
	if err != nil {
		return err
	}

	if len(requests) == 0 {
		return fmt.Errorf("no pending requests")
	}

	if reply != "agree" {
		for _, req := range requests {
			_ = s.groupJoinRequest.UpdateStatus(ctx, req.ID, model.GroupJoinRequestRejected)
		}
		return nil
	}

	for _, req := range requests {
		groupMember := &model.GroupMember{
			GroupID: gid,
			UserID:  req.UserID,
			Role:    1,
		}

		if err := s.groupMemberRepo.AddMember(ctx, groupMember); err != nil {
			return fmt.Errorf("add member failed: %w", err)
		}

		if err := s.groupJoinRequest.UpdateStatus(ctx, req.ID, model.GroupJoinRequestAccepted); err != nil {
			return err
		}
	}

	return nil
}

func (s *GroupService) IsGroupMember(ctx context.Context, gid, uid string) (bool, error) {
	return s.groupMemberRepo.IsMember(ctx, gid, uid)
}

func (s *GroupService) GetPendingGroupRequests(ctx context.Context, gid string) ([]*model.GroupJoinRequest, error) {
	return s.groupJoinRequest.GetPendingRequests(ctx, gid)
}

func (s *GroupService) GetGroupsByUser(ctx context.Context, uid string) ([]*model.Group, error) {
	return s.groupRepo.GetGroupsByUserID(ctx, uid)
}

func (s *GroupService) RemoveMember(ctx context.Context, gid, ownerID, memberID string) error {
	group, err := s.groupRepo.GetByID(ctx, gid)
	if err != nil {
		return err
	}

	if group.OwnerID != ownerID {
		return model.ErrNotOwner
	}

	isMember, err := s.groupMemberRepo.IsMember(ctx, gid, memberID)
	if err != nil {
		return err
	}
	if !isMember {
		return model.ErrNotMember
	}

	return s.groupMemberRepo.RemoveMember(ctx, gid, memberID)
}

func (s *GroupService) GetGroupMembers(ctx context.Context, gid string) ([]*model.User, error) {
	return s.groupMemberRepo.GetMembers(ctx, gid)
}

func (s *GroupService) GetMemberRole(ctx context.Context, gid, uid string) (int16, error) {
	return s.groupMemberRepo.GetRole(ctx, gid, uid)
}

func (s *GroupService) UpdateMemberRole(ctx context.Context, gid, ownerID, uid string, role int16) error {
	group, err := s.groupRepo.GetByID(ctx, gid)
	if err != nil {
		return err
	}

	if group.OwnerID != ownerID {
		return model.ErrNotOwner
	}

	return s.groupMemberRepo.UpdateRole(ctx, gid, uid, role)
}
