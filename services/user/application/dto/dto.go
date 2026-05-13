package dto

type RegisterRequest struct {
	Tele     string `json:"tele" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Tele     string `json:"tele" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Token  string `json:"token"`
}

type UserInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Tele      string `json:"tele"`
	AvatarURL string `json:"avatar_url"`
	Status    int    `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type UpdateUserRequest struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type AddFriendRequest struct {
	FriendID string `json:"friend_id" binding:"required"`
	Reason   string `json:"reason"`
}

type ReplyFriendRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Accept    bool   `json:"accept"`
}

type FriendRequestInfo struct {
	ID        string    `json:"id"`
	FromUID   string    `json:"from_uid"`
	ToUID     string    `json:"to_uid"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	CreatedAt int64     `json:"created_at"`
	FromUser  *UserInfo `json:"from_user,omitempty"`
}
