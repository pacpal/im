// Package service 提供 group 服务的业务逻辑实现。
package service

import (
	"IM/pkg/id"
	"IM/services/group/domain/entity"
	"IM/services/group/domain/event"
	"IM/services/group/domain/repository"
	"context"
	"errors"
	"time"
)

var (
	ErrGroupNotFound      = errors.New("group not found")
	ErrAlreadyMember      = errors.New("already a member")
	ErrNotMember          = errors.New("not a member")
	ErrNotOwner           = errors.New("not group owner")
	ErrNotAdmin           = errors.New("not admin")
	ErrRequestNotFound    = errors.New("request not found")
	ErrRequestExists      = errors.New("request already exists")
	ErrCannotRemoveOwner  = errors.New("cannot remove owner")
	ErrCannotLeaveAsOwner = errors.New("owner cannot leave group, transfer ownership first")
)

// GroupService 提供群组的创建、更新、加入、处理请求等业务逻辑。
type GroupService struct {
	groupRepo            repository.GroupRepository
	groupMemberRepo      repository.GroupMemberRepository
	groupJoinRequestRepo repository.GroupJoinRequestRepository
	idGenerator          *id.Generator
	eventPublisher       *event.EventPublisher
}

// NewGroupService 创建 GroupService 实例。
func NewGroupService(
	groupRepo repository.GroupRepository,
	groupMemberRepo repository.GroupMemberRepository,
	groupJoinRequestRepo repository.GroupJoinRequestRepository,
	idGenerator *id.Generator,
	eventPublisher *event.EventPublisher,
) *GroupService {
	return &GroupService{
		groupRepo:            groupRepo,
		groupMemberRepo:      groupMemberRepo,
		groupJoinRequestRepo: groupJoinRequestRepo,
		idGenerator:          idGenerator,
		eventPublisher:       eventPublisher,
	}
}

func (s *GroupService) CreateGroup(ctx context.Context, ownerID, name, description string) (*entity.Group, error) {
	groupID := s.idGenerator.Generate()
	group := entity.NewGroup(groupID, name, ownerID, description)

	if err := s.groupRepo.Create(ctx, group); err != nil {
		return nil, err
	}

	owner := entity.NewGroupMember(groupID, ownerID, entity.MemberRoleOwner)
	if err := s.groupMemberRepo.Create(ctx, owner); err != nil {
		return nil, err
	}

	s.eventPublisher.Publish(&event.GroupCreatedEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "group.created",
			OccurredAt:  time.Now(),
			AggregateID: groupID,
		},
		GroupID: groupID,
		Name:    name,
		OwnerID: ownerID,
	})

	return group, nil
}

func (s *GroupService) GetGroup(ctx context.Context, groupID string) (*entity.Group, error) {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return nil, ErrGroupNotFound
	}
	return group, nil
}

func (s *GroupService) UpdateGroup(ctx context.Context, groupID, ownerID, name, description, imageURL string) error {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}

	if !group.IsOwner(ownerID) {
		return ErrNotOwner
	}

	group.Update(name, description, imageURL)
	if err := s.groupRepo.Update(ctx, group); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.GroupUpdatedEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "group.updated",
			OccurredAt:  time.Now(),
			AggregateID: groupID,
		},
		GroupID: groupID,
		Name:    group.Name,
	})

	return nil
}

func (s *GroupService) DeleteGroup(ctx context.Context, groupID, ownerID string) error {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}

	if !group.IsOwner(ownerID) {
		return ErrNotOwner
	}

	if err := s.groupRepo.Delete(ctx, groupID); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.GroupDeletedEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "group.deleted",
			OccurredAt:  time.Now(),
			AggregateID: groupID,
		},
		GroupID: groupID,
		OwnerID: ownerID,
	})

	return nil
}

func (s *GroupService) JoinGroup(ctx context.Context, userID, groupID, reason string) error {
	isMember, err := s.groupMemberRepo.Exists(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return ErrAlreadyMember
	}

	requestExists, err := s.groupJoinRequestRepo.Exists(ctx, userID, groupID)
	if err != nil {
		return err
	}
	if requestExists {
		return ErrRequestExists
	}

	requestID := s.idGenerator.Generate()
	request := entity.NewGroupJoinRequest(requestID, userID, groupID, reason)

	if err := s.groupJoinRequestRepo.Create(ctx, request); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.JoinRequestCreatedEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "group_join_request.created",
			OccurredAt:  time.Now(),
			AggregateID: requestID,
		},
		RequestID: requestID,
		UserID:    userID,
		GroupID:   groupID,
	})

	return nil
}

func (s *GroupService) LeaveGroup(ctx context.Context, userID, groupID string) error {
	member, err := s.groupMemberRepo.GetByGroupAndUserID(ctx, groupID, userID)
	if err != nil {
		return ErrNotMember
	}

	if member.IsOwner() {
		return ErrCannotLeaveAsOwner
	}

	if err := s.groupMemberRepo.Delete(ctx, groupID, userID); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.MemberLeftEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "group.member_left",
			OccurredAt:  time.Now(),
			AggregateID: groupID,
		},
		GroupID: groupID,
		UserID:  userID,
	})

	return nil
}

func (s *GroupService) AcceptJoinRequest(ctx context.Context, requestID, ownerID string) error {
	request, err := s.groupJoinRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		return ErrRequestNotFound
	}

	group, err := s.groupRepo.GetByID(ctx, request.GroupID)
	if err != nil {
		return ErrGroupNotFound
	}

	if !group.IsOwner(ownerID) {
		return ErrNotOwner
	}

	if !request.IsPending() {
		return ErrRequestNotFound
	}

	member := entity.NewGroupMember(request.GroupID, request.UserID, entity.MemberRoleMember)
	if err := s.groupMemberRepo.Create(ctx, member); err != nil {
		return err
	}

	request.Accept()
	if err := s.groupJoinRequestRepo.UpdateStatus(ctx, requestID, entity.RequestStatusAccepted); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.JoinRequestAcceptedEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "group_join_request.accepted",
			OccurredAt:  time.Now(),
			AggregateID: requestID,
		},
		RequestID: requestID,
		UserID:    request.UserID,
		GroupID:   request.GroupID,
	})

	return nil
}

func (s *GroupService) RejectJoinRequest(ctx context.Context, requestID, ownerID string) error {
	request, err := s.groupJoinRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		return ErrRequestNotFound
	}

	group, err := s.groupRepo.GetByID(ctx, request.GroupID)
	if err != nil {
		return ErrGroupNotFound
	}

	if !group.IsOwner(ownerID) {
		return ErrNotOwner
	}

	if !request.IsPending() {
		return ErrRequestNotFound
	}

	request.Reject()
	if err := s.groupJoinRequestRepo.UpdateStatus(ctx, requestID, entity.RequestStatusRejected); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.JoinRequestRejectedEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "group_join_request.rejected",
			OccurredAt:  time.Now(),
			AggregateID: requestID,
		},
		RequestID: requestID,
		UserID:    request.UserID,
		GroupID:   request.GroupID,
	})

	return nil
}

func (s *GroupService) GetMembers(ctx context.Context, groupID string) ([]*entity.GroupMember, error) {
	return s.groupMemberRepo.GetByGroupID(ctx, groupID)
}

func (s *GroupService) RemoveMember(ctx context.Context, groupID, ownerID, memberID string) error {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}

	if !group.IsOwner(ownerID) {
		return ErrNotOwner
	}

	member, err := s.groupMemberRepo.GetByGroupAndUserID(ctx, groupID, memberID)
	if err != nil {
		return ErrNotMember
	}

	if member.IsOwner() {
		return ErrCannotRemoveOwner
	}

	if err := s.groupMemberRepo.Delete(ctx, groupID, memberID); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.MemberRemovedEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "group.member_removed",
			OccurredAt:  time.Now(),
			AggregateID: groupID,
		},
		GroupID: groupID,
		UserID:  memberID,
		OwnerID: ownerID,
	})

	return nil
}

func (s *GroupService) GetUserGroups(ctx context.Context, userID string) ([]*entity.Group, error) {
	return s.groupRepo.GetByUserID(ctx, userID)
}

func (s *GroupService) GetPendingRequests(ctx context.Context, groupID string) ([]*entity.GroupJoinRequest, error) {
	return s.groupJoinRequestRepo.GetPendingByGroupID(ctx, groupID)
}

func (s *GroupService) TransferOwner(ctx context.Context, groupID, ownerID, newOwnerID string) error {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}

	if !group.IsOwner(ownerID) {
		return ErrNotOwner
	}

	newOwner, err := s.groupMemberRepo.GetByGroupAndUserID(ctx, groupID, newOwnerID)
	if err != nil {
		return ErrNotMember
	}

	oldOwner, err := s.groupMemberRepo.GetByGroupAndUserID(ctx, groupID, ownerID)
	if err != nil {
		return ErrNotMember
	}

	oldOwner.SetRole(entity.MemberRoleMember)
	newOwner.SetRole(entity.MemberRoleOwner)

	group.OwnerID = newOwnerID
	group.UpdatedAt = time.Now()

	if err := s.groupRepo.Update(ctx, group); err != nil {
		return err
	}

	s.eventPublisher.Publish(&event.OwnerTransferredEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "group.owner_transferred",
			OccurredAt:  time.Now(),
			AggregateID: groupID,
		},
		GroupID:    groupID,
		OldOwnerID: ownerID,
		NewOwnerID: newOwnerID,
	})

	return nil
}
