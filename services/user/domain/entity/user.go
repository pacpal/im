// Package entity 定义 user 域的实体：User、Friendship、FriendRequest 及其行为。
package entity

import "time"

// User 表示系统中的用户实体。
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

// UserStatus 表示用户在线状态的枚举。
type UserStatus int

const (
	UserStatusOffline UserStatus = 0
	UserStatusOnline  UserStatus = 1
	UserStatusBusy    UserStatus = 2
)

// NewUser 创建一个新的 User 实例并设置默认在线状态与时间戳。
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

// UpdateProfile 更新用户的显示名与头像并更新时间戳。
func (u *User) UpdateProfile(name, avatarURL string) {
	if name != "" {
		u.Name = name
	}
	if avatarURL != "" {
		u.AvatarURL = avatarURL
	}
	u.UpdatedAt = time.Now()
}

// SetStatus 设置用户状态并更新时间戳。
func (u *User) SetStatus(status UserStatus) {
	u.Status = status
	u.UpdatedAt = time.Now()
}

// IsValid 做一些基本合法性校验（仅示例）。
func (u *User) IsValid() bool {
	return u.Name != "" && u.Tele != ""
}

// Friendship 表示一条双向好友关系中的单向记录。
type Friendship struct {
	UserID    string
	FriendID  string
	Status    FriendshipStatus
	CreatedAt time.Time
}

// FriendshipStatus 表示好友关系的状态。
type FriendshipStatus int

const (
	FriendshipStatusActive  FriendshipStatus = 1
	FriendshipStatusBlocked FriendshipStatus = 2
)

// NewFriendship 创建一条新的 Friendship 记录，默认为激活状态。
func NewFriendship(userID, friendID string) *Friendship {
	return &Friendship{
		UserID:    userID,
		FriendID:  friendID,
		Status:    FriendshipStatusActive,
		CreatedAt: time.Now(),
	}
}

// FriendRequest 表示一次好友申请。
type FriendRequest struct {
	ID        string
	FromUID   string
	ToUID     string
	Reason    string
	Status    FriendRequestStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// FriendRequestStatus 定义好友申请的状态字符串。
type FriendRequestStatus string

const (
	FriendRequestStatusPending  FriendRequestStatus = "pending"
	FriendRequestStatusAccepted FriendRequestStatus = "accepted"
	FriendRequestStatusRejected FriendRequestStatus = "rejected"
)

// NewFriendRequest 创建一个初始状态为 pending 的好友申请。
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

// Accept 将申请状态标记为已接受并更新时间戳。
func (r *FriendRequest) Accept() {
	r.Status = FriendRequestStatusAccepted
	r.UpdatedAt = time.Now()
}

// Reject 将申请状态标记为已拒绝并更新时间戳。
func (r *FriendRequest) Reject() {
	r.Status = FriendRequestStatusRejected
	r.UpdatedAt = time.Now()
}

// IsPending 判断申请是否处于 pending 状态。
func (r *FriendRequest) IsPending() bool {
	return r.Status == FriendRequestStatusPending
}
