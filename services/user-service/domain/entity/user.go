package entity

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Tele      string    `json:"tele"`
	Password  string    `json:"-"`
	AvatarURL string    `json:"avatar_url"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) IsValid() bool {
	return u.Name != "" && u.Tele != ""
}

type Friendship struct {
	UserID    string    `json:"user_id"`
	FriendID  string    `json:"friend_id"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	FriendshipStatusActive = 1
	FriendshipStatusBlocked = 2
)

type FriendRequest struct {
	ID        string    `json:"id"`
	FromUID   string    `json:"from_uid"`
	ToUID     string    `json:"to_uid"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	FriendRequestStatusPending  = "pending"
	FriendRequestStatusAccepted = "accepted"
	FriendRequestStatusRejected = "rejected"
)