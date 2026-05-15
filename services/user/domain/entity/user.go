package entity

import "time"

type User struct {
	ID        string
	Name      string
	Tele      string
	Password  string
	AvatarURL string
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserStatus int

const (
	UserStatusOffline UserStatus = 0
	UserStatusOnline  UserStatus = 1
	UserStatusBusy    UserStatus = 2
)

func NewUser(id, name, tele, password string) *User {
	return &User{
		ID:        id,
		Name:      name,
		Tele:      tele,
		Password:  password,
		Status:    UserStatusOnline,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) UpdateProfile(name, avatarURL string) {
	if name != "" {
		u.Name = name
	}
	if avatarURL != "" {
		u.AvatarURL = avatarURL
	}
	u.UpdatedAt = time.Now()
}

func (u *User) SetStatus(status UserStatus) {
	u.Status = status
	u.UpdatedAt = time.Now()
}

func (u *User) IsValid() bool {
	return u.Name != "" && u.Tele != ""
}

type Friendship struct {
	UserID    string
	FriendID  string
	Status    FriendshipStatus
	CreatedAt time.Time
}

type FriendshipStatus int

const (
	FriendshipStatusActive  FriendshipStatus = 1
	FriendshipStatusBlocked FriendshipStatus = 2
)

func NewFriendship(userID, friendID string) *Friendship {
	return &Friendship{
		UserID:    userID,
		FriendID:  friendID,
		Status:    FriendshipStatusActive,
		CreatedAt: time.Now(),
	}
}

type FriendRequest struct {
	ID        string
	FromUID   string
	ToUID     string
	Reason    string
	Status    FriendRequestStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FriendRequestStatus string

const (
	FriendRequestStatusPending  FriendRequestStatus = "pending"
	FriendRequestStatusAccepted FriendRequestStatus = "accepted"
	FriendRequestStatusRejected FriendRequestStatus = "rejected"
)

func NewFriendRequest(id, fromUID, toUID, reason string) *FriendRequest {
	return &FriendRequest{
		ID:        id,
		FromUID:   fromUID,
		ToUID:     toUID,
		Reason:    reason,
		Status:    FriendRequestStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (r *FriendRequest) Accept() {
	r.Status = FriendRequestStatusAccepted
	r.UpdatedAt = time.Now()
}

func (r *FriendRequest) Reject() {
	r.Status = FriendRequestStatusRejected
	r.UpdatedAt = time.Now()
}

func (r *FriendRequest) IsPending() bool {
	return r.Status == FriendRequestStatusPending
}
