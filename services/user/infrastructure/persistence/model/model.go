package model

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;size:64"`
	Name      string    `gorm:"size:100;not null;index"`
	Tele      string    `gorm:"size:20;not null;uniqueIndex"`
	Password  string    `gorm:"size:255;not null"`
	AvatarURL string    `gorm:"column:avatar_url;size:255"`
	Status    int       `gorm:"default:1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

type Friendship struct {
	UserID    string    `gorm:"primaryKey;type:varchar(64);column:user_id"`
	FriendID  string    `gorm:"primaryKey;type:varchar(64);column:friend_id"`
	Status    int       `gorm:"type:smallint;default:1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Friendship) TableName() string {
	return "friendships"
}

type FriendRequest struct {
	ID        string    `gorm:"primaryKey;size:64"`
	FromUID   string    `gorm:"column:from_uid;size:64;not null;index"`
	ToUID     string    `gorm:"column:to_uid;size:64;not null;index"`
	Reason    string    `gorm:"type:text"`
	Status    string    `gorm:"size:20;not null;default:'pending';index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (FriendRequest) TableName() string {
	return "friend_requests"
}
