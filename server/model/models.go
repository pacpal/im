// Package model store struct msg,user,group
package model

import "errors"

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
	SdID    string `json:"send_id"`
	RcID    string `json:"receive_id"`
	Content string `json:"content"`
	Type    string `json:"type"`
	Time    int64  `json:"time"`
}

type User struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Friends  map[string]bool `json:"friends,omitempty"`
	Groups   map[string]bool `json:"groups,omitempty"`
	Tele     string
	Password string
}

type Group struct {
	ID        string
	Name      string
	MemberIDs map[string]bool // key userid
	// MemberCache map[string]*UserCache
	Description string
	OwnerID     string
	Type        string
	ImageURL    string
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
