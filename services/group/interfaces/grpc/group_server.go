package grpc

import (
	"IM/api/gen/group"
	"IM/services/group/application/service"
	"context"
)

type GroupServer struct {
	grouppb.UnimplementedGroupServiceServer
	groupSvc *service.GroupService
}

func NewGroupServer(groupSvc *service.GroupService) *GroupServer {
	return &GroupServer{
		groupSvc: groupSvc,
	}
}

func (s *GroupServer) CreateGroup(ctx context.Context, req *grouppb.CreateGroupRequest) (*grouppb.CreateGroupResponse, error) {
	g, err := s.groupSvc.CreateGroup(ctx, req.OwnerId, req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	return &grouppb.CreateGroupResponse{
		GroupId: g.ID,
		Name:    g.Name,
	}, nil
}

func (s *GroupServer) GetGroup(ctx context.Context, req *grouppb.GetGroupRequest) (*grouppb.GroupInfo, error) {
	g, err := s.groupSvc.GetGroup(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}

	return &grouppb.GroupInfo{
		Id:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		OwnerId:     g.OwnerID,
		Type:        string(g.Type),
		ImageUrl:    g.ImageURL,
		CreatedAt:   g.CreatedAt.Unix(),
	}, nil
}

func (s *GroupServer) UpdateGroup(ctx context.Context, req *grouppb.UpdateGroupRequest) (*grouppb.CommonResponse, error) {
	err := s.groupSvc.UpdateGroup(ctx, req.GroupId, req.OwnerId, req.Name, req.Description, req.ImageUrl)
	if err != nil {
		return &grouppb.CommonResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &grouppb.CommonResponse{
		Success: true,
		Message: "group updated",
	}, nil
}

func (s *GroupServer) DeleteGroup(ctx context.Context, req *grouppb.DeleteGroupRequest) (*grouppb.CommonResponse, error) {
	err := s.groupSvc.DeleteGroup(ctx, req.GroupId, req.OwnerId)
	if err != nil {
		return &grouppb.CommonResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &grouppb.CommonResponse{
		Success: true,
		Message: "group deleted",
	}, nil
}

func (s *GroupServer) JoinGroup(ctx context.Context, req *grouppb.JoinGroupRequest) (*grouppb.CommonResponse, error) {
	err := s.groupSvc.JoinGroup(ctx, req.UserId, req.GroupId, req.Reason)
	if err != nil {
		return &grouppb.CommonResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &grouppb.CommonResponse{
		Success: true,
		Message: "join request sent",
	}, nil
}

func (s *GroupServer) LeaveGroup(ctx context.Context, req *grouppb.LeaveGroupRequest) (*grouppb.CommonResponse, error) {
	err := s.groupSvc.LeaveGroup(ctx, req.UserId, req.GroupId)
	if err != nil {
		return &grouppb.CommonResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &grouppb.CommonResponse{
		Success: true,
		Message: "left group",
	}, nil
}

func (s *GroupServer) ReplyGroupJoinRequest(ctx context.Context, req *grouppb.ReplyGroupJoinRequest) (*grouppb.CommonResponse, error) {
	var err error
	if req.Accept {
		err = s.groupSvc.AcceptJoinRequest(ctx, req.RequestId, req.OwnerId)
	} else {
		err = s.groupSvc.RejectJoinRequest(ctx, req.RequestId, req.OwnerId)
	}

	if err != nil {
		return &grouppb.CommonResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &grouppb.CommonResponse{
		Success: true,
		Message: "request processed",
	}, nil
}

func (s *GroupServer) GetMembers(ctx context.Context, req *grouppb.GetMembersRequest) (*grouppb.GetMembersResponse, error) {
	members, err := s.groupSvc.GetMembers(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}

	pbMembers := make([]*grouppb.MemberInfo, len(members))
	for i, m := range members {
		pbMembers[i] = &grouppb.MemberInfo{
			UserId:   m.UserID,
			Role:     int32(m.Role),
			JoinedAt: m.JoinedAt.Unix(),
			Nickname: m.Nickname,
		}
	}

	return &grouppb.GetMembersResponse{
		Members: pbMembers,
	}, nil
}

func (s *GroupServer) RemoveMember(ctx context.Context, req *grouppb.RemoveMemberRequest) (*grouppb.CommonResponse, error) {
	err := s.groupSvc.RemoveMember(ctx, req.GroupId, req.OwnerId, req.MemberId)
	if err != nil {
		return &grouppb.CommonResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &grouppb.CommonResponse{
		Success: true,
		Message: "member removed",
	}, nil
}

func (s *GroupServer) GetUserGroups(ctx context.Context, req *grouppb.GetUserGroupsRequest) (*grouppb.GetUserGroupsResponse, error) {
	groups, err := s.groupSvc.GetUserGroups(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	pbGroups := make([]*grouppb.GroupInfo, len(groups))
	for i, g := range groups {
		pbGroups[i] = &grouppb.GroupInfo{
			Id:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			OwnerId:     g.OwnerID,
			Type:        string(g.Type),
			ImageUrl:    g.ImageURL,
			CreatedAt:   g.CreatedAt.Unix(),
		}
	}

	return &grouppb.GetUserGroupsResponse{
		Groups: pbGroups,
	}, nil
}

func (s *GroupServer) GetPendingGroupRequests(ctx context.Context, req *grouppb.GetPendingGroupRequestsRequest) (*grouppb.GetPendingGroupRequestsResponse, error) {
	requests, err := s.groupSvc.GetPendingRequests(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}

	pbRequests := make([]*grouppb.GroupJoinRequestInfo, len(requests))
	for i, r := range requests {
		pbRequests[i] = &grouppb.GroupJoinRequestInfo{
			Id:        r.ID,
			UserId:    r.UserID,
			GroupId:   r.GroupID,
			Reason:    r.Reason,
			Status:    string(r.Status),
			CreatedAt: r.CreatedAt.Unix(),
		}
	}

	return &grouppb.GetPendingGroupRequestsResponse{
		Requests: pbRequests,
	}, nil
}

func (s *GroupServer) TransferOwner(ctx context.Context, req *grouppb.TransferOwnerRequest) (*grouppb.CommonResponse, error) {
	err := s.groupSvc.TransferOwner(ctx, req.GroupId, req.OwnerId, req.NewOwnerId)
	if err != nil {
		return &grouppb.CommonResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &grouppb.CommonResponse{
		Success: true,
		Message: "ownership transferred",
	}, nil
}
