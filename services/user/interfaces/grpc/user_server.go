package grpc

import (
	"IM/api/gen/user"
	"IM/services/user/application/service"
	"context"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	userSvc *service.UserService
}

func NewUserServer(userSvc *service.UserService) *UserServer {
	return &UserServer{
		userSvc: userSvc,
	}
}

func (s *UserServer) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	user, err := s.userSvc.Register(ctx, req.Tele, req.Name, req.Password)
	if err != nil {
		return &userpb.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &userpb.RegisterResponse{
		UserId: user.ID,
		Name:   user.Name,
	}, nil
}

func (s *UserServer) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	user, token, err := s.userSvc.Login(ctx, req.Tele, req.Password)
	if err != nil {
		return &userpb.LoginResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &userpb.LoginResponse{
		UserId: user.ID,
		Name:   user.Name,
		Token:  token,
	}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserInfo, error) {
	user, err := s.userSvc.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &userpb.UserInfo{
		Id:        user.ID,
		Name:      user.Name,
		Tele:      user.Tele,
		AvatarUrl: user.AvatarURL,
		Status:    int32(user.Status),
		CreatedAt: user.CreatedAt.Unix(),
	}, nil
}

func (s *UserServer) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	users, err := s.userSvc.GetUsersByIDs(ctx, req.UserIds)
	if err != nil {
		return nil, err
	}

	pbUsers := make([]*userpb.UserInfo, len(users))
	for i, u := range users {
		pbUsers[i] = &userpb.UserInfo{
			Id:        u.ID,
			Name:      u.Name,
			Tele:      u.Tele,
			AvatarUrl: u.AvatarURL,
			Status:    int32(u.Status),
			CreatedAt: u.CreatedAt.Unix(),
		}
	}

	return &userpb.GetUsersResponse{
		Users: pbUsers,
	}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.CommonResponse, error) {
	user, err := s.userSvc.GetUserByID(ctx, req.UserId)
	if err != nil {
		return &userpb.CommonResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	user.UpdateProfile(req.Name, req.AvatarUrl)
	if err := s.userSvc.UpdateUser(ctx, user); err != nil {
		return &userpb.CommonResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &userpb.CommonResponse{
		Success: true,
		Message: "user updated",
	}, nil
}

func (s *UserServer) GetFriends(ctx context.Context, req *userpb.GetFriendsRequest) (*userpb.GetFriendsResponse, error) {
	friends, err := s.userSvc.GetFriends(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	pbFriends := make([]*userpb.UserInfo, len(friends))
	for i, f := range friends {
		pbFriends[i] = &userpb.UserInfo{
			Id:        f.ID,
			Name:      f.Name,
			Tele:      f.Tele,
			AvatarUrl: f.AvatarURL,
			Status:    int32(f.Status),
		}
	}

	return &userpb.GetFriendsResponse{
		Friends: pbFriends,
	}, nil
}

func (s *UserServer) AddFriend(ctx context.Context, req *userpb.AddFriendRequest) (*userpb.CommonResponse, error) {
	err := s.userSvc.AddFriend(ctx, req.UserId, req.TargetId, req.Reason)
	if err != nil {
		return &userpb.CommonResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &userpb.CommonResponse{
		Success: true,
		Message: "friend request sent",
	}, nil
}

func (s *UserServer) RemoveFriend(ctx context.Context, req *userpb.RemoveFriendRequest) (*userpb.CommonResponse, error) {
	err := s.userSvc.RemoveFriend(ctx, req.UserId, req.TargetId)
	if err != nil {
		return &userpb.CommonResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &userpb.CommonResponse{
		Success: true,
		Message: "friend removed",
	}, nil
}

func (s *UserServer) ReplyFriendRequest(ctx context.Context, req *userpb.ReplyFriendRequest) (*userpb.CommonResponse, error) {
	var err error
	if req.Accept {
		err = s.userSvc.AcceptFriendRequest(ctx, req.RequestId, req.UserId)
	} else {
		err = s.userSvc.RejectFriendRequest(ctx, req.RequestId, req.UserId)
	}

	if err != nil {
		return &userpb.CommonResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &userpb.CommonResponse{
		Success: true,
		Message: "friend request processed",
	}, nil
}

func (s *UserServer) CheckFriendship(ctx context.Context, req *userpb.CheckFriendshipRequest) (*userpb.CheckFriendshipResponse, error) {
	isFriend, err := s.userSvc.CheckFriendship(ctx, req.UserId1, req.UserId2)
	if err != nil {
		return nil, err
	}

	return &userpb.CheckFriendshipResponse{
		IsFriend: isFriend,
	}, nil
}

func (s *UserServer) GetPendingFriendRequests(ctx context.Context, req *userpb.GetPendingFriendRequestsRequest) (*userpb.GetPendingFriendRequestsResponse, error) {
	requests, err := s.userSvc.GetPendingFriendRequests(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	pbRequests := make([]*userpb.FriendRequestInfo, len(requests))
	for i, r := range requests {
		pbRequests[i] = &userpb.FriendRequestInfo{
			Id:        r.ID,
			FromUid:   r.FromUID,
			ToUid:     r.ToUID,
			Reason:    r.Reason,
			Status:    string(r.Status),
			CreatedAt: r.CreatedAt.Unix(),
		}
	}

	return &userpb.GetPendingFriendRequestsResponse{
		Requests: pbRequests,
	}, nil
}
