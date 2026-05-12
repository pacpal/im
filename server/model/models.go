package model

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrNotFriend       = errors.New("not friend")
	ErrAlreadyFriend   = errors.New("already friend")
	ErrGroupNotFound   = errors.New("group not found")
	ErrAlreadyMember   = errors.New("already member")
	ErrNotMember       = errors.New("not member")
	ErrRequestNotFound = errors.New("request not found")
	ErrInvalidRequest  = errors.New("invalid request")
	ErrMessageNotFound = errors.New("message not found")
	ErrNotOwner        = errors.New("not group owner")
)

type Message struct {
	ID        string    `json:"msg_id" gorm:"primaryKey;size:64"`
	SdID      string    `json:"send_id" gorm:"column:send_id;size:64;not null;index"`
	RcID      string    `json:"receive_id" gorm:"column:receive_id;size:64;not null;index"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Type      string    `json:"type" gorm:"size:20;not null;index"`
	Time      int64     `json:"time" gorm:"not null;index"`
	IsRead    bool      `json:"is_read" gorm:"default:false;index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (m *Message) TableName() string {
	return "messages"
}

type User struct {
	ID        string    `json:"id" gorm:"primaryKey;size:64"`
	Name      string    `json:"name" gorm:"size:100;not null;index"`
	Tele      string    `json:"tele" gorm:"size:20;not null;uniqueIndex"`
	Password  string    `json:"-" gorm:"size:255;not null"`
	AvatarURL string    `json:"avatar_url" gorm:"column:avatar_url;size:255"`
	Status    int       `json:"status" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Friends []*User  `json:"friends,omitempty" gorm:"many2many:friendships;foreignKey:ID;joinForeignKey:user_id;references:ID;joinReferences:friend_id;"`
	Groups  []*Group `json:"groups,omitempty" gorm:"many2many:group_members;foreignKey:ID;joinForeignKey:user_id;references:ID;joinReferences:group_id;"`
}

func (u *User) TableName() string {
	return "users"
}

type Friendship struct {
	UserID    string    `gorm:"primaryKey;type:varchar(64);column:user_id"`
	FriendID  string    `gorm:"primaryKey;type:varchar(64);column:friend_id"`
	Status    int       `gorm:"type:smallint;default:1;comment:'1:好友 2:拉黑'"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
}

func (f *Friendship) TableName() string {
	return "friendships"
}

type Group struct {
	ID          string    `json:"id" gorm:"primaryKey;size:64"`
	Name        string    `json:"name" gorm:"size:100;not null;index"`
	Description string    `json:"description" gorm:"type:text"`
	OwnerID     string    `json:"owner_id" gorm:"column:owner_id;size:64;not null;index"`
	Type        string    `json:"type" gorm:"size:20;not null;default:'normal'"`
	ImageURL    string    `json:"image_url" gorm:"column:image_url;size:255"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Owner   *User   `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	Members []*User `json:"members,omitempty" gorm:"many2many:group_members;foreignKey:ID;joinForeignKey:group_id;references:ID;joinReferences:user_id;"`
}

func (g *Group) TableName() string {
	return "groups"
}

func NewGroup(id, name, ownerID string) *Group {
	return &Group{
		ID:      id,
		Name:    name,
		OwnerID: ownerID,
		Type:    "normal",
		Members: make([]*User, 0),
	}
}

type GroupMember struct {
	GroupID  string    `gorm:"primaryKey;type:varchar(64);column:group_id"`
	UserID   string    `gorm:"primaryKey;type:varchar(64);column:user_id"`
	Role     int16     `gorm:"column:role;type:smallint;default:1;comment:'1:普通成员 2:管理员 3:群主'"`
	JoinedAt time.Time `gorm:"autoCreateTime;column:joined_at"`
	Nickname string    `gorm:"column:nickname;size:100"`
}

func (gm *GroupMember) TableName() string {
	return "group_members"
}

type FriendRequest struct {
	ID        string    `json:"id" gorm:"primaryKey;size:64"`
	FromUID   string    `json:"from_uid" gorm:"column:from_uid;size:64;not null;index"`
	ToUID     string    `json:"to_uid" gorm:"column:to_uid;size:64;not null;index"`
	Reason    string    `json:"reason" gorm:"type:text"`
	Status    string    `json:"status" gorm:"size:20;not null;default:'pending';index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (f *FriendRequest) TableName() string {
	return "friend_requests"
}

const (
	FriendRequestPending  = "pending"
	FriendRequestAccepted = "accepted"
	FriendRequestRejected = "rejected"
)

type GroupJoinRequest struct {
	ID        string    `json:"id" gorm:"primaryKey;size:64"`
	UserID    string    `json:"user_id" gorm:"column:user_id;size:64;not null;index"`
	GroupID   string    `json:"group_id" gorm:"column:group_id;size:64;not null;index"`
	Reason    string    `json:"reason" gorm:"type:text"`
	Status    string    `json:"status" gorm:"size:20;not null;default:'pending';index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (g *GroupJoinRequest) TableName() string {
	return "group_join_requests"
}

const (
	GroupJoinRequestPending  = "pending"
	GroupJoinRequestAccepted = "accepted"
	GroupJoinRequestRejected = "rejected"
)
