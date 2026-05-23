// Package entity 定义 group 域的实体及其行为。
package entity

import "time"

// Group 表示群组实体。
type Group struct {
	ID          string
	Name        string
	Description string
	OwnerID     string
	Type        GroupType
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GroupType 表示群组类型。
type GroupType string

const (
	GroupTypeNormal  GroupType = "normal"
	GroupTypePrivate GroupType = "private"
	GroupTypePublic  GroupType = "public"
)

// NewGroup 创建一个新的 Group 实例。
func NewGroup(id, name, ownerID, description string) *Group {
	return &Group{
		ID:          id,
		Name:        name,
		OwnerID:     ownerID,
		Description: description,
		Type:        GroupTypeNormal,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// Update 更新群组信息并更新时间戳。
func (g *Group) Update(name, description, imageURL string) {
	if name != "" {
		g.Name = name
	}
	if description != "" {
		g.Description = description
	}
	if imageURL != "" {
		g.ImageURL = imageURL
	}
	g.UpdatedAt = time.Now()
}

// IsOwner 判断给定用户是否为群主。
func (g *Group) IsOwner(userID string) bool {
	return g.OwnerID == userID
}

// GroupMember 表示群成员。
type GroupMember struct {
	GroupID  string
	UserID   string
	Role     MemberRole
	Nickname string
	JoinedAt time.Time
}

// MemberRole 表示群成员角色的枚举。
type MemberRole int

const (
	MemberRoleMember MemberRole = 1
	MemberRoleAdmin  MemberRole = 2
	MemberRoleOwner  MemberRole = 3
)

// NewGroupMember 创建新的群成员记录。
func NewGroupMember(groupID, userID string, role MemberRole) *GroupMember {
	return &GroupMember{
		GroupID:  groupID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}
}

// IsOwner 判断成员是否为群主。
func (m *GroupMember) IsOwner() bool {
	return m.Role == MemberRoleOwner
}

// IsAdmin 判断成员是否至少为管理员。
func (m *GroupMember) IsAdmin() bool {
	return m.Role >= MemberRoleAdmin
}

// SetRole 设置成员角色。
func (m *GroupMember) SetRole(role MemberRole) {
	m.Role = role
}

// GroupJoinRequest 表示用户加入群组的请求。
type GroupJoinRequest struct {
	ID        string
	UserID    string
	GroupID   string
	Reason    string
	Status    RequestStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RequestStatus 表示加入请求的状态字符串类型。
type RequestStatus string

const (
	RequestStatusPending  RequestStatus = "pending"
	RequestStatusAccepted RequestStatus = "accepted"
	RequestStatusRejected RequestStatus = "rejected"
)

// NewGroupJoinRequest 创建一个 join 请求并设置为 pending。
func NewGroupJoinRequest(id, userID, groupID, reason string) *GroupJoinRequest {
	return &GroupJoinRequest{
		ID:        id,
		UserID:    userID,
		GroupID:   groupID,
		Reason:    reason,
		Status:    RequestStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Accept 接受请求并更新时间戳。
func (r *GroupJoinRequest) Accept() {
	r.Status = RequestStatusAccepted
	r.UpdatedAt = time.Now()
}

// Reject 拒绝请求并更新时间戳。
func (r *GroupJoinRequest) Reject() {
	r.Status = RequestStatusRejected
	r.UpdatedAt = time.Now()
}

// IsPending 判断请求是否为 pending。
func (r *GroupJoinRequest) IsPending() bool {
	return r.Status == RequestStatusPending
}
