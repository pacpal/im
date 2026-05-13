package service

import (
	"IM/services/group-service/domain/entity"
	"IM/services/group-service/infrastructure/repository"
	"context"
	"errors"
	"time"
)

var (
	ErrGroupNotFound      = errors.New("group not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrAlreadyMember     = errors.New("already a member")
	ErrNotMember         = errors.New("not a member")
	ErrNotOwner          = errors.New("not the group owner")
	ErrNotAdmin          = errors.New("not an admin")
	ErrRequestNotFound   = errors.New("request not found")
	ErrRequestExists     = errors.New("join request already exists")
	ErrInvalidGroupName  = errors.New("invalid group name")
)

type GroupApplicationService struct {
	groupRepo        repository.GroupRepository
	groupMemberRepo  repository.GroupMemberRepository
	groupJoinRequest repository.GroupJoinRequestRepository
	userRepo         repository.UserRepository
}

func NewGroupApplicationService(
	groupRepo repository.GroupRepository,
	groupMemberRepo repository.GroupMemberRepository,
	groupJoinRequest repository.GroupJoinRequestRepository,
	userRepo repository.UserRepository,
) *GroupApplicationService {
	return &GroupApplicationService{
		groupRepo:        groupRepo,
		groupMemberRepo:  groupMemberRepo,
		groupJoinRequest: groupJoinRequest,
		userRepo:         userRepo,
	}
}

func (s *GroupApplicationService) CreateGroup(ctx context.Context, ownerID, name, description string) (*entity.Group, error) {
	if len(name) < 2 || len(name) > 100 {
		return nil, ErrInvalidGroupName
	}

	exists, err := s.userRepo.Exists(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	group := &entity.Group{
		ID:          generateGroupID(),
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		Type:        "normal",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.groupRepo.Create(ctx, group); err != nil {
		return nil, err
	}

	member := &entity.GroupMember{
		GroupID:  group.ID,
		UserID:   ownerID,
		Role:     entity.GroupRoleOwner,
		JoinedAt: time.Now(),
	}

	if err := s.groupMemberRepo.AddMember(ctx, member); err != nil {
		return nil, err
	}

	return group, nil
}

func (s *GroupApplicationService) GetGroup(ctx context.Context, groupID string) (*entity.Group, error) {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return nil, ErrGroupNotFound
	}
	return group, nil
}

func (s *GroupApplicationService) GetGroupMembers(ctx context.Context, groupID string) ([]*entity.User, error) {
	return s.groupMemberRepo.GetMembers(ctx, groupID)
}

func (s *GroupApplicationService) GetUserGroups(ctx context.Context, userID string) ([]*entity.Group, error) {
	return s.groupRepo.GetGroupsByUserID(ctx, userID)
}

func (s *GroupApplicationService) JoinGroup(ctx context.Context, userID, groupID, reason string) error {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}
	if group == nil {
		return ErrGroupNotFound
	}

	isMember, err := s.groupMemberRepo.IsMember(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return ErrAlreadyMember
	}

	exists, err := s.groupJoinRequest.Exists(ctx, userID, groupID)
	if err != nil {
		return err
	}
	if exists {
		return ErrRequestExists
	}

	request := &entity.GroupJoinRequest{
		ID:        generateRequestID("greq"),
		UserID:    userID,
		GroupID:   groupID,
		Reason:    reason,
		Status:    entity.GroupJoinRequestStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.groupJoinRequest.Create(ctx, request)
}

func (s *GroupApplicationService) AcceptJoinRequest(ctx context.Context, requestID, approverID string) error {
	req, err := s.groupJoinRequest.GetByID(ctx, requestID)
	if err != nil {
		return ErrRequestNotFound
	}

	role, err := s.groupMemberRepo.GetRole(ctx, req.GroupID, approverID)
	if err != nil {
		return err
	}
	if role != entity.GroupRoleOwner && role != entity.GroupRoleAdmin {
		return ErrNotAdmin
	}

	member := &entity.GroupMember{
		GroupID:  req.GroupID,
		UserID:   req.UserID,
		Role:     entity.GroupRoleNormal,
		JoinedAt: time.Now(),
	}

	if err := s.groupMemberRepo.AddMember(ctx, member); err != nil {
		return err
	}

	return s.groupJoinRequest.UpdateStatus(ctx, requestID, entity.GroupJoinRequestStatusAccepted)
}

func (s *GroupApplicationService) RejectJoinRequest(ctx context.Context, requestID, rejecterID string) error {
	req, err := s.groupJoinRequest.GetByID(ctx, requestID)
	if err != nil {
		return ErrRequestNotFound
	}

	role, err := s.groupMemberRepo.GetRole(ctx, req.GroupID, rejecterID)
	if err != nil {
		return err
	}
	if role != entity.GroupRoleOwner && role != entity.GroupRoleAdmin {
		return ErrNotAdmin
	}

	return s.groupJoinRequest.UpdateStatus(ctx, requestID, entity.GroupJoinRequestStatusRejected)
}

func (s *GroupApplicationService) LeaveGroup(ctx context.Context, userID, groupID string) error {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}

	if group.OwnerID == userID {
		return ErrNotOwner
	}

	return s.groupMemberRepo.RemoveMember(ctx, groupID, userID)
}

func (s *GroupApplicationService) RemoveMember(ctx context.Context, operatorID, targetID, groupID string) error {
	role, err := s.groupMemberRepo.GetRole(ctx, groupID, operatorID)
	if err != nil {
		return err
	}
	if role != entity.GroupRoleOwner && role != entity.GroupRoleAdmin {
		return ErrNotAdmin
	}

	group, _ := s.groupRepo.GetByID(ctx, groupID)
	if group != nil && group.OwnerID == targetID {
		return ErrNotOwner
	}

	return s.groupMemberRepo.RemoveMember(ctx, groupID, targetID)
}

func (s *GroupApplicationService) UpdateMemberRole(ctx context.Context, operatorID, targetID, groupID string, newRole int16) error {
	role, err := s.groupMemberRepo.GetRole(ctx, groupID, operatorID)
	if err != nil {
		return err
	}
	if role != entity.GroupRoleOwner {
		return ErrNotOwner
	}

	return s.groupMemberRepo.UpdateRole(ctx, groupID, targetID, newRole)
}

func (s *GroupApplicationService) GetPendingJoinRequests(ctx context.Context, groupID string) ([]*entity.GroupJoinRequest, error) {
	return s.groupJoinRequest.GetPendingRequests(ctx, groupID)
}

func generateGroupID() string {
	return "group_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

func generateRequestID(prefix string) string {
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