package dto

import "IM/services/user-service/domain/entity"

type RegisterRequest struct {
	Tele     string `json:"tele"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Tele     string `json:"tele"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Tele      string `json:"tele"`
	AvatarURL string `json:"avatar_url"`
	Status    int    `json:"status"`
}

type AddFriendRequest struct {
	ToUID   string `json:"to_uid"`
	Reason  string `json:"reason"`
}

type FriendRequestResponse struct {
	ID        string `json:"id"`
	FromUID   string `json:"from_uid"`
	ToUID     string `json:"to_uid"`
	Reason    string `json:"reason"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type FriendRequestActionRequest struct {
	RequestID string `json:"request_id"`
}

func ToUserResponse(user *entity.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Tele:      user.Tele,
		AvatarURL: user.AvatarURL,
		Status:    user.Status,
	}
}

func ToUserResponseList(users []*entity.User) []*UserResponse {
	result := make([]*UserResponse, len(users))
	for i, u := range users {
		result[i] = ToUserResponse(u)
	}
	return result
}

func ToFriendRequestResponse(req *entity.FriendRequest) *FriendRequestResponse {
	if req == nil {
		return nil
	}
	return &FriendRequestResponse{
		ID:        req.ID,
		FromUID:   req.FromUID,
		ToUID:     req.ToUID,
		Reason:    req.Reason,
		Status:    req.Status,
		CreatedAt: req.CreatedAt.Unix(),
	}
}

func ToFriendRequestResponseList(reqs []*entity.FriendRequest) []*FriendRequestResponse {
	result := make([]*FriendRequestResponse, len(reqs))
	for i, r := range reqs {
		result[i] = ToFriendRequestResponse(r)
	}
	return result
}