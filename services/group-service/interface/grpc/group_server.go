package grpc

import (
	"IM/services/group-service/application/service"
	pb "IM/api/gen/group"
	"context"
)

type GroupGrpcServer struct {
	pb.UnimplementedGroupServiceServer
	groupSvc *service.GroupApplicationService
}

func NewGroupGrpcServer(groupSvc *service.GroupApplicationService) *GroupGrpcServer {
	return &GroupGrpcServer{
		groupSvc: groupSvc,
	}
}

func (s *GroupGrpcServer) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	group, err := s.groupSvc.CreateGroup(ctx, req.OwnerId, req.Name, req.Description)
	if err != nil {
		return &pb.CreateGroupResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.CreateGroupResponse{
		Success: true,
		Message: "group created",
		Group: &pb.Group{
			Id:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			OwnerId:     group.OwnerID,
			Type:        group.Type,
		},
	}, nil
}

func (s *GroupGrpcServer) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GetGroupResponse, error) {
	group, err := s.groupSvc.GetGroup(ctx, req.GroupId)
	if err != nil {
		return &pb.GetGroupResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.GetGroupResponse{
		Success: true,
		Group: &pb.Group{
			Id:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			OwnerId:     group.OwnerID,
			Type:        group.Type,
		},
	}, nil
}

func (s *GroupGrpcServer) GetGroupMembers(ctx context.Context, req *pb.GetGroupMembersRequest) (*pb.GetGroupMembersResponse, error) {
	members, err := s.groupSvc.GetGroupMembers(ctx, req.GroupId)
	if err != nil {
		return &pb.GetGroupMembersResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	pbMembers := make([]*pb.GroupMember, len(members))
	for i, m := range members {
		pbMembers[i] = &pb.GroupMember{
			UserId: m.ID,
			Name:   m.Name,
			Role:   0,
		}
	}

	return &pb.GetGroupMembersResponse{
		Success: true,
		Members: pbMembers,
	}, nil
}

func (s *GroupGrpcServer) JoinGroup(ctx context.Context, req *pb.JoinGroupRequest) (*pb.JoinGroupResponse, error) {
	err := s.groupSvc.JoinGroup(ctx, req.UserId, req.GroupId, req.Reason)
	if err != nil {
		return &pb.JoinGroupResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.JoinGroupResponse{
		Success: true,
		Message: "join request sent",
	}, nil
}

func (s *GroupGrpcServer) AcceptJoinRequest(ctx context.Context, req *pb.AcceptJoinRequestRequest) (*pb.AcceptJoinRequestResponse, error) {
	err := s.groupSvc.AcceptJoinRequest(ctx, req.RequestId, req.ApproverId)
	if err != nil {
		return &pb.AcceptJoinRequestResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.AcceptJoinRequestResponse{
		Success: true,
		Message: "join request accepted",
	}, nil
}

func (s *GroupGrpcServer) LeaveGroup(ctx context.Context, req *pb.LeaveGroupRequest) (*pb.LeaveGroupResponse, error) {
	err := s.groupSvc.LeaveGroup(ctx, req.UserId, req.GroupId)
	if err != nil {
		return &pb.LeaveGroupResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.LeaveGroupResponse{
		Success: true,
		Message: "left group",
	}, nil
}