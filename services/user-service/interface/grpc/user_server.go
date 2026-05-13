package grpc

import (
	"IM/services/user-service/application/service"
	pb "IM/api/gen/user"
	"context"
)

type UserGrpcServer struct {
	pb.UnimplementedUserServiceServer
	userSvc *service.UserApplicationService
}

func NewUserGrpcServer(userSvc *service.UserApplicationService) *UserGrpcServer {
	return &UserGrpcServer{
		userSvc: userSvc,
	}
}

func (s *UserGrpcServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := s.userSvc.Register(ctx, req.Tele, req.Name, req.Password)
	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "register success",
		User: &pb.User{
			Id:     user.ID,
			Name:   user.Name,
			Tele:   user.Tele,
			Status: int32(user.Status),
		},
	}, nil
}

func (s *UserGrpcServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.userSvc.Login(ctx, req.Tele, req.Password)
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.LoginResponse{
		Success: true,
		Message: "login success",
		User: &pb.User{
			Id:     user.ID,
			Name:   user.Name,
			Tele:   user.Tele,
			Status: int32(user.Status),
		},
	}, nil
}

func (s *UserGrpcServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.userSvc.GetUserByID(ctx, req.UserId)
	if err != nil {
		return &pb.GetUserResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.GetUserResponse{
		Success: true,
		User: &pb.User{
			Id:     user.ID,
			Name:   user.Name,
			Tele:   user.Tele,
			Status: int32(user.Status),
		},
	}, nil
}

func (s *UserGrpcServer) GetFriends(ctx context.Context, req *pb.GetFriendsRequest) (*pb.GetFriendsResponse, error) {
	friends, err := s.userSvc.GetFriends(ctx, req.UserId)
	if err != nil {
		return &pb.GetFriendsResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	pbFriends := make([]*pb.User, len(friends))
	for i, f := range friends {
		pbFriends[i] = &pb.User{
			Id:     f.ID,
			Name:   f.Name,
			Tele:   f.Tele,
			Status: int32(f.Status),
		}
	}

	return &pb.GetFriendsResponse{
		Success: true,
		Users:   pbFriends,
	}, nil
}

func (s *UserGrpcServer) AddFriend(ctx context.Context, req *pb.AddFriendRequest) (*pb.AddFriendResponse, error) {
	err := s.userSvc.AddFriend(ctx, req.UserId, req.FriendId, req.Reason)
	if err != nil {
		return &pb.AddFriendResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.AddFriendResponse{
		Success: true,
		Message: "friend request sent",
	}, nil
}

func (s *UserGrpcServer) AcceptFriendRequest(ctx context.Context, req *pb.AcceptFriendRequestRequest) (*pb.AcceptFriendRequestResponse, error) {
	err := s.userSvc.AcceptFriendRequest(ctx, req.RequestId, req.UserId)
	if err != nil {
		return &pb.AcceptFriendRequestResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.AcceptFriendRequestResponse{
		Success: true,
		Message: "friend request accepted",
	}, nil
}

func (s *UserGrpcServer) GetPendingFriendRequests(ctx context.Context, req *pb.GetPendingFriendRequestsRequest) (*pb.GetPendingFriendRequestsResponse, error) {
	requests, err := s.userSvc.GetPendingFriendRequests(ctx, req.UserId)
	if err != nil {
		return &pb.GetPendingFriendRequestsResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	pbRequests := make([]*pb.FriendRequest, len(requests))
	for i, r := range requests {
		pbRequests[i] = &pb.FriendRequest{
			Id:      r.ID,
			FromUid: r.FromUID,
			ToUid:   r.ToUID,
			Reason:  r.Reason,
			Status:  r.Status,
		}
	}

	return &pb.GetPendingFriendRequestsResponse{
		Success:  true,
		Requests: pbRequests,
	}, nil
}