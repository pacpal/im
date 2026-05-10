// Package model store struct msg,user,group
package model

import (
	"errors"
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
)

type Message struct {
	ID      string `json:"msg_id" gorm:"column:id;primaryKey;size:64"`
	SdID    string `json:"send_id" gorm:"column:send_id;size:64;not null"`
	RcID    string `json:"receive_id" gorm:"column:receive_id;size:64;not null"`
	Content string `json:"content" gorm:"column:content;type:text;not null"`
	Type    string `json:"type" gorm:"column:type;size:20;not null"`
	Time    int64  `json:"time" gorm:"column:time"`
}

type User struct {
	ID       string `json:"id" gorm:"column:id;primaryKey;size:64"`
	Name     string `json:"name" gorm:"column:name;size:100;not null"`
	Tele     string `json:"tele" gorm:"column:tele;size:20;not null"`
	Password string `gorm:"column:password;size:255;not null"`
	//change
	Friends map[string]bool `json:"friends,omitempty"`
	Groups  map[string]bool `json:"groups,omitempty"`
}

type Group struct {
	ID          string          `json:"id" gorm:"size:64;primaryKey"`
	Name        string          `json:"name" gorm:"size:100;not null"`
	Description string          `json:"description" gorm:"type:text"`
	OwnerID     string          `json:"ownerid" gorm:"size:64;not null"`
	Type        string          `json:"type" gorm:"size:20;not null"`
	ImageURL    string          `json:"imgurl"`
	MemberIDs   map[string]bool // key userid
	// MemberCache map[string]*UserCache
}
type GroupMember struct {
	GroupID string `gorm:"`
}

func NewGroup(id, name, ownerID string) *Group {
	g := &Group{
		ID:        id,
		Name:      name,
		OwnerID:   ownerID,
		MemberIDs: make(map[string]bool),
		Type:      "normal",
	}
	g.MemberIDs[ownerID] = true
	return g
}
func (u *User) AddFriend(uid string) error {
	if u.Friends[uid] {
		return ErrAlreadyFriend
	}
	u.Friends[uid] = true
	return nil
}

func (u *User) RemoveFriend(uid string) error {
	if !u.Friends[uid] {
		return ErrNotFriend
	}
	delete(u.Friends, uid)
	return nil
}

func (u *User) JoinGroup(gid string) error {
	if u.Groups[gid] {
		return ErrAlreadyMember
	}
	u.Groups[gid] = true
	return nil
}

func (u *User) LeaveGroup(gid string) error {
	if !u.Groups[gid] {
		return ErrNotMember
	}
	delete(u.Groups, gid)
	return nil
}
func (g *Group) AddMember(uid string) error {
	if g.MemberIDs[uid] {
		return ErrAlreadyMember
	}
	g.MemberIDs[uid] = true
	return nil
}

func (g *Group) RemoveMember(uid string) error {
	if !g.MemberIDs[uid] {
		return ErrNotMember
	}
	delete(g.MemberIDs, uid)
	return nil
}

func (g *Group) IsMember(uid string) bool {
	return g.MemberIDs[uid]
}
