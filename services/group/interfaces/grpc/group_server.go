package grpc

import (
	"IM/api/gen/group"
	common "IM/api/gen/common"
	"IM/services/group/application/service"
	"context"
)

type GroupServer struct {
	group.UnimplementedGroupServiceServer
	groupSvc *service.GroupService
}

func NewGroupServer(groupSvc *service.GroupService) *GroupServer {
	return &GroupServer{
		groupSvc: groupSvc,
	}
}

func (s *GroupServer) CreateGroup(ctx context.Context, req *group.CreateGroupRequest) (*group.CreateGroupResponse, error) {
	g, err := s.groupSvc.CreateGroup(ctx, req.OwnerId, req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	return &group.CreateGroupResponse{
		GroupId: g.ID,
		Name:    g.Name,
	}, nil
}

func (s *GroupServer) GetGroup(ctx context.Context, req *group.GetGroupRequest) (*group.GroupInfo, error) {
	g, err := s.groupSvc.GetGroup(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}

	return &group.GroupInfo{
		Id:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		OwnerId:     g.OwnerID,
		Type:        string(g.Type),
		ImageUrl:    g.ImageURL,
		CreatedAt:   g.CreatedAt.Unix(),
	}, nil
}

func (s *GroupServer) JoinGroup(ctx context.Context, req *group.JoinGroupRequest) (*common.Response, error) {
	err := s.groupSvc.JoinGroup(ctx, req.UserId, req.GroupId, req.Reason)
	if err != nil {
		return &common.Response{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &common.Response{
		Success: true,
		Message: "join request sent",
	}, nil
}

func (s *GroupServer) LeaveGroup(ctx context.Context, req *group.LeaveGroupRequest) (*common.Response, error) {
	err := s.groupSvc.LeaveGroup(ctx, req.UserId, req.GroupId)
	if err != nil {
		return &common.Response{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &common.Response{
		Success: true,
		Message: "left group",
	}, nil
}

func (s *GroupServer) ReplyGroupJoin(ctx context.Context, req *group.ReplyGroupJoinRequest) (*common.Response, error) {
	var err error
	if req.GetAccept() {
		err = s.groupSvc.AcceptJoinRequest(ctx, req.GetRequestId(), req.GetOwnerId())
	} else {
		err = s.groupSvc.RejectJoinRequest(ctx, req.GetRequestId(), req.GetOwnerId())
	}

	if err != nil {
		return &common.Response{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &common.Response{
		Success: true,
		Message: "request processed",
	}, nil
}

func (s *GroupServer) GetMembers(ctx context.Context, req *group.GetMembersRequest) (*group.GetMembersResponse, error) {
	members, err := s.groupSvc.GetMembers(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}

	pbMembers := make([]*group.MemberInfo, len(members))
	for i, m := range members {
		pbMembers[i] = &group.MemberInfo{
			UserId:   m.UserID,
			Role:     int32(m.Role),
			JoinedAt: m.JoinedAt.Unix(),
		}
	}

	return &group.GetMembersResponse{
		Members: pbMembers,
	}, nil
}

func (s *GroupServer) RemoveMember(ctx context.Context, req *group.RemoveMemberRequest) (*common.Response, error) {
	err := s.groupSvc.RemoveMember(ctx, req.GroupId, req.OwnerId, req.MemberId)
	if err != nil {
		return &common.Response{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &common.Response{
		Success: true,
		Message: "member removed",
	}, nil
}

func (s *GroupServer) GetUserGroups(ctx context.Context, req *group.GetUserGroupsRequest) (*group.GetUserGroupsResponse, error) {
	groups, err := s.groupSvc.GetUserGroups(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	pbGroups := make([]*group.GroupInfo, len(groups))
	for i, g := range groups {
		pbGroups[i] = &group.GroupInfo{
			Id:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			OwnerId:     g.OwnerID,
			Type:        string(g.Type),
			ImageUrl:    g.ImageURL,
			CreatedAt:   g.CreatedAt.Unix(),
		}
	}

	return &group.GetUserGroupsResponse{
		Groups: pbGroups,
	}, nil
}

func (s *GroupServer) GetPendingGroupRequests(ctx context.Context, req *group.GetPendingGroupRequestsRequest) (*group.GetPendingGroupRequestsResponse, error) {
	requests, err := s.groupSvc.GetPendingRequests(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}

	pbRequests := make([]*group.GroupJoinRequestInfo, len(requests))
	for i, r := range requests {
		pbRequests[i] = &group.GroupJoinRequestInfo{
			Id:        r.ID,
			UserId:    r.UserID,
			GroupId:   r.GroupID,
			Reason:    r.Reason,
			Status:    string(r.Status),
			CreatedAt: r.CreatedAt.Unix(),
		}
	}

	return &group.GetPendingGroupRequestsResponse{
		Requests: pbRequests,
	}, nil
}
