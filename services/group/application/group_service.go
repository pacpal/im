// Package service 提供 group 服务的业务逻辑实现。
package service

import (
	"IM/pkg/id"
	"IM/pkg/logger"
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

func (s *GroupService) CreateGroup(ctx context.Context, ownerID, name, description string) (res *entity.Group, err error) {
	done := logger.StartStep("GroupService.CreateGroup", "owner", ownerID, "name", name)
	defer func() { done(err) }()

	groupID := s.idGenerator.Generate()
	res = entity.NewGroup(groupID, name, ownerID, description)

	if err = s.groupRepo.Create(ctx, res); err != nil {
		return nil, err
	}

	owner := entity.NewGroupMember(groupID, ownerID, entity.MemberRoleOwner)
	if err = s.groupMemberRepo.Create(ctx, owner); err != nil {
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

	logger.Infow("CreateGroup: created", "component", "group_service", "group_id", groupID)
	return res, nil
}

func (s *GroupService) GetGroup(ctx context.Context, groupID string) (res *entity.Group, err error) {
	done := logger.StartStep("GroupService.GetGroup", "group_id", groupID)
	defer func() { done(err) }()

	res, err = s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		err = ErrGroupNotFound
		return
	}
	logger.Infow("GetGroup: found", "component", "group_service", "group_id", res.ID)
	return
}

func (s *GroupService) UpdateGroup(ctx context.Context, groupID, ownerID, name, description, imageURL string) (err error) {
	done := logger.StartStep("GroupService.UpdateGroup", "group_id", groupID, "owner", ownerID)
	defer func() { done(err) }()

	var group *entity.Group
	group, err = s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		err = ErrGroupNotFound
		return
	}

	if !group.IsOwner(ownerID) {
		err = ErrNotOwner
		return
	}

	group.Update(name, description, imageURL)
	if err = s.groupRepo.Update(ctx, group); err != nil {
		return
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

	logger.Infow("UpdateGroup: updated", "component", "group_service", "group_id", groupID)
	return
}

func (s *GroupService) DeleteGroup(ctx context.Context, groupID, ownerID string) (err error) {
	done := logger.StartStep("GroupService.DeleteGroup", "group_id", groupID, "owner", ownerID)
	defer func() { done(err) }()

	var group *entity.Group
	group, err = s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		err = ErrGroupNotFound
		return
	}

	if !group.IsOwner(ownerID) {
		err = ErrNotOwner
		return
	}

	if err = s.groupRepo.Delete(ctx, groupID); err != nil {
		return
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

	logger.Infow("DeleteGroup: deleted", "component", "group_service", "group_id", groupID)
	return
}

func (s *GroupService) JoinGroup(ctx context.Context, userID, groupID, reason string) (err error) {
	done := logger.StartStep("GroupService.JoinGroup", "user", userID, "group", groupID)
	defer func() { done(err) }()

	var isMember bool
	isMember, err = s.groupMemberRepo.Exists(ctx, groupID, userID)
	if err != nil {
		return
	}
	if isMember {
		err = ErrAlreadyMember
		return
	}

	var requestExists bool
	requestExists, err = s.groupJoinRequestRepo.Exists(ctx, userID, groupID)
	if err != nil {
		return
	}
	if requestExists {
		err = ErrRequestExists
		return
	}

	requestID := s.idGenerator.Generate()
	request := entity.NewGroupJoinRequest(requestID, userID, groupID, reason)

	if err = s.groupJoinRequestRepo.Create(ctx, request); err != nil {
		return
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

	logger.Infow("JoinGroup: join request created", "component", "group_service", "request_id", requestID)
	return
}

func (s *GroupService) LeaveGroup(ctx context.Context, userID, groupID string) (err error) {
	done := logger.StartStep("GroupService.LeaveGroup", "user", userID, "group", groupID)
	defer func() { done(err) }()

	var member *entity.GroupMember
	member, err = s.groupMemberRepo.GetByGroupAndUserID(ctx, groupID, userID)
	if err != nil {
		err = ErrNotMember
		return
	}

	if member.IsOwner() {
		err = ErrCannotLeaveAsOwner
		return
	}

	if err = s.groupMemberRepo.Delete(ctx, groupID, userID); err != nil {
		return
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

	logger.Infow("LeaveGroup: user left group", "component", "group_service", "group_id", groupID, "user", userID)
	return
}

func (s *GroupService) AcceptJoinRequest(ctx context.Context, requestID, ownerID string) (err error) {
	done := logger.StartStep("GroupService.AcceptJoinRequest", "request_id", requestID, "owner", ownerID)
	defer func() { done(err) }()

	var request *entity.GroupJoinRequest
	request, err = s.groupJoinRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		err = ErrRequestNotFound
		return
	}

	var group *entity.Group
	group, err = s.groupRepo.GetByID(ctx, request.GroupID)
	if err != nil {
		err = ErrGroupNotFound
		return
	}

	if !group.IsOwner(ownerID) {
		err = ErrNotOwner
		return
	}

	if !request.IsPending() {
		err = ErrRequestNotFound
		return
	}

	member := entity.NewGroupMember(request.GroupID, request.UserID, entity.MemberRoleMember)
	if err = s.groupMemberRepo.Create(ctx, member); err != nil {
		return
	}

	request.Accept()
	if err = s.groupJoinRequestRepo.UpdateStatus(ctx, requestID, entity.RequestStatusAccepted); err != nil {
		return
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

	logger.Infow("AcceptJoinRequest: accepted", "component", "group_service", "request_id", requestID)
	return
}

func (s *GroupService) RejectJoinRequest(ctx context.Context, requestID, ownerID string) (err error) {
	done := logger.StartStep("GroupService.RejectJoinRequest", "request_id", requestID, "owner", ownerID)
	defer func() { done(err) }()

	var request *entity.GroupJoinRequest
	request, err = s.groupJoinRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		err = ErrRequestNotFound
		return
	}

	var group *entity.Group
	group, err = s.groupRepo.GetByID(ctx, request.GroupID)
	if err != nil {
		err = ErrGroupNotFound
		return
	}

	if !group.IsOwner(ownerID) {
		err = ErrNotOwner
		return
	}

	if !request.IsPending() {
		err = ErrRequestNotFound
		return
	}

	request.Reject()
	if err = s.groupJoinRequestRepo.UpdateStatus(ctx, requestID, entity.RequestStatusRejected); err != nil {
		return
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

	logger.Infow("RejectJoinRequest: rejected", "component", "group_service", "request_id", requestID)
	return
}

func (s *GroupService) GetMembers(ctx context.Context, groupID string) (res []*entity.GroupMember, err error) {
	done := logger.StartStep("GroupService.GetMembers", "group_id", groupID)
	defer func() { done(err) }()

	res, err = s.groupMemberRepo.GetByGroupID(ctx, groupID)
	if err == nil {
		logger.Infow("GetMembers: retrieved", "component", "group_service", "group_id", groupID, "count", len(res))
	}
	return
}

func (s *GroupService) RemoveMember(ctx context.Context, groupID, adminID, memberID string) (err error) {
	done := logger.StartStep("GroupService.RemoveMember", "group_id", groupID, "admin", adminID, "member", memberID)
	defer func() { done(err) }()

	_, err = s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		err = ErrGroupNotFound
		return
	}

	var member *entity.GroupMember
	member, err = s.groupMemberRepo.GetByGroupAndUserID(ctx, groupID, adminID)
	if err != nil {
		err = ErrNotMember
		return
	}

	if member.IsAdmin() {
		err = ErrCannotRemoveOwner
		return
	}

	if err = s.groupMemberRepo.Delete(ctx, groupID, memberID); err != nil {
		return
	}

	s.eventPublisher.Publish(&event.MemberRemovedEvent{
		BaseEvent: event.BaseEvent{
			EventType:   "group.member_removed",
			OccurredAt:  time.Now(),
			AggregateID: groupID,
		},
		GroupID: groupID,
		UserID:  memberID,
		AdminID: adminID,
	})

	logger.Infow("RemoveMember: removed", "component", "group_service", "group_id", groupID, "member", memberID)
	return
}

func (s *GroupService) GetUserGroups(ctx context.Context, userID string) (res []*entity.Group, err error) {
	done := logger.StartStep("GroupService.GetUserGroups", "user_id", userID)
	defer func() { done(err) }()

	res, err = s.groupRepo.GetByUserID(ctx, userID)
	if err == nil {
		logger.Infow("GetUserGroups: retrieved", "component", "group_service", "user_id", userID, "count", len(res))
	}
	return
}

func (s *GroupService) GetPendingRequests(ctx context.Context, groupID string) (res []*entity.GroupJoinRequest, err error) {
	done := logger.StartStep("GroupService.GetPendingRequests", "group_id", groupID)
	defer func() { done(err) }()

	res, err = s.groupJoinRequestRepo.GetPendingByGroupID(ctx, groupID)
	if err == nil {
		logger.Infow("GetPendingRequests: retrieved", "component", "group_service", "group_id", groupID, "count", len(res))
	}
	return
}

func (s *GroupService) TransferOwner(ctx context.Context, groupID, ownerID, newOwnerID string) (err error) {
	done := logger.StartStep("GroupService.TransferOwner", "group_id", groupID, "old_owner", ownerID, "new_owner", newOwnerID)
	defer func() { done(err) }()

	var group *entity.Group
	group, err = s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		err = ErrGroupNotFound
		return
	}

	if !group.IsOwner(ownerID) {
		err = ErrNotOwner
		return
	}

	var newOwner *entity.GroupMember
	newOwner, err = s.groupMemberRepo.GetByGroupAndUserID(ctx, groupID, newOwnerID)
	if err != nil {
		err = ErrNotMember
		return
	}

	var oldOwner *entity.GroupMember
	oldOwner, err = s.groupMemberRepo.GetByGroupAndUserID(ctx, groupID, ownerID)
	if err != nil {
		err = ErrNotMember
		return
	}

	oldOwner.SetRole(entity.MemberRoleMember)
	newOwner.SetRole(entity.MemberRoleOwner)

	group.OwnerID = newOwnerID
	group.UpdatedAt = time.Now()

	if err = s.groupRepo.Update(ctx, group); err != nil {
		return
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

	logger.Infow("TransferOwner: transferred", "component", "group_service", "group_id", groupID, "old_owner", ownerID, "new_owner", newOwnerID)
	return
}
