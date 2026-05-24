package grpc

import (
	"IM/api/gen/common"
	"IM/api/gen/user"
	"IM/services/user/application/service"
	"context"
)

type UserServer struct {
	user.UnimplementedUserServiceServer
	userSvc *service.UserService
}

func NewUserServer(userSvc *service.UserService) *UserServer {
	return &UserServer{
		userSvc: userSvc,
	}
}

func (s *UserServer) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	u, err := s.userSvc.Register(ctx, req.Tele, req.Name, req.Password)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResponse{
		UserId: u.ID,
		Name:   u.Name,
	}, nil
}

func (s *UserServer) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	u, token, err := s.userSvc.Login(ctx, req.Tele, req.Id, req.Password)
	if err != nil {
		return nil, err
	}

	return &user.LoginResponse{
		UserId: u.ID,
		Name:   u.Name,
		Token:  token,
	}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.UserInfo, error) {
	u, err := s.userSvc.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &user.UserInfo{
		Id:        u.ID,
		Name:      u.Name,
		Tele:      u.Tele,
		AvatarUrl: u.AvatarURL,
		Status:    int32(u.Status),
		CreatedAt: u.CreatedAt.Unix(),
	}, nil
}

func (s *UserServer) GetFriends(ctx context.Context, req *user.GetFriendsRequest) (*user.GetFriendsResponse, error) {
	friends, err := s.userSvc.GetFriends(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	pbFriends := make([]*user.UserInfo, len(friends))
	for i, f := range friends {
		pbFriends[i] = &user.UserInfo{
			Id:        f.ID,
			Name:      f.Name,
			Tele:      f.Tele,
			AvatarUrl: f.AvatarURL,
			Status:    int32(f.Status),
		}
	}

	return &user.GetFriendsResponse{
		Friends: pbFriends,
	}, nil
}

func (s *UserServer) AddFriend(ctx context.Context, req *user.AddFriendRequest) (*common.Response, error) {
	err := s.userSvc.AddFriend(ctx, req.UserId, req.TargetId, req.Reason)
	if err != nil {
		return &common.Response{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &common.Response{
		Success: true,
		Message: "friend request sent",
	}, nil
}

func (s *UserServer) RemoveFriend(ctx context.Context, req *user.RemoveFriendRequest) (*common.Response, error) {
	err := s.userSvc.RemoveFriend(ctx, req.UserId, req.TargetId)
	if err != nil {
		return &common.Response{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &common.Response{
		Success: true,
		Message: "friend removed",
	}, nil
}

func (s *UserServer) ReplyFriend(ctx context.Context, req *user.ReplyFriendRequest) (*common.Response, error) {
	var err error
	if req.GetAccept() {
		err = s.userSvc.AcceptFriendRequest(ctx, req.RequestId, req.UserId)
	} else {
		err = s.userSvc.RejectFriendRequest(ctx, req.RequestId, req.UserId)
	}

	if err != nil {
		return &common.Response{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &common.Response{
		Success: true,
		Message: "friend request processed",
	}, nil
}

func (s *UserServer) CheckFriendship(ctx context.Context, req *user.CheckFriendshipRequest) (*user.CheckFriendshipResponse, error) {
	isFriend, err := s.userSvc.CheckFriendship(ctx, req.UserId1, req.UserId2)
	if err != nil {
		return nil, err
	}

	return &user.CheckFriendshipResponse{
		IsFriend: isFriend,
	}, nil
}

func (s *UserServer) GetPendingFriendRequests(ctx context.Context, req *user.GetPendingFriendRequestsRequest) (*user.GetPendingFriendRequestsResponse, error) {
	requests, err := s.userSvc.GetPendingFriendRequests(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	pbRequests := make([]*user.FriendRequestInfo, len(requests))
	for i, r := range requests {
		pbRequests[i] = &user.FriendRequestInfo{
			Id:        r.ID,
			FromUid:   r.FromUID,
			ToUid:     r.ToUID,
			Reason:    r.Reason,
			Status:    string(r.Status),
			CreatedAt: r.CreatedAt.Unix(),
		}
	}

	return &user.GetPendingFriendRequestsResponse{
		Requests: pbRequests,
	}, nil
}
