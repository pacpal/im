// Package model 定义 group 服务使用的 GORM 模型结构体。
package model

import "time"

// Group 为 groups 表的映射模型。
type Group struct {
	ID          string    `gorm:"primaryKey;size:64"`
	Name        string    `gorm:"size:100;not null;index"`
	Description string    `gorm:"type:text"`
	OwnerID     string    `gorm:"column:owner_id;size:64;not null;index"`
	Type        string    `gorm:"size:20;not null;default:'normal'"`
	ImageURL    string    `gorm:"column:image_url;size:255"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (Group) TableName() string {
	return "groups"
}

type GroupMember struct {
	GroupID  string    `gorm:"primaryKey;type:varchar(64);column:group_id"`
	UserID   string    `gorm:"primaryKey;type:varchar(64);column:user_id"`
	Role     int       `gorm:"column:role;type:smallint;default:1"`
	Nickname string    `gorm:"column:nickname;size:100"`
	JoinedAt time.Time `gorm:"autoCreateTime"`
}

func (GroupMember) TableName() string {
	return "group_members"
}

type GroupJoinRequest struct {
	ID        string    `gorm:"primaryKey;size:64"`
	UserID    string    `gorm:"column:user_id;size:64;not null;index"`
	GroupID   string    `gorm:"column:group_id;size:64;not null;index"`
	Reason    string    `gorm:"type:text"`
	Status    string    `gorm:"size:20;not null;default:'pending';index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (GroupJoinRequest) TableName() string {
	return "group_join_requests"
}
