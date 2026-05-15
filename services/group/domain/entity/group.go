package entity

import "time"

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

type GroupType string

const (
	GroupTypeNormal  GroupType = "normal"
	GroupTypePrivate GroupType = "private"
	GroupTypePublic  GroupType = "public"
)

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

func (g *Group) IsOwner(userID string) bool {
	return g.OwnerID == userID
}

type GroupMember struct {
	GroupID  string
	UserID   string
	Role     MemberRole
	Nickname string
	JoinedAt time.Time
}

type MemberRole int

const (
	MemberRoleMember MemberRole = 1
	MemberRoleAdmin  MemberRole = 2
	MemberRoleOwner  MemberRole = 3
)

func NewGroupMember(groupID, userID string, role MemberRole) *GroupMember {
	return &GroupMember{
		GroupID:  groupID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}
}

func (m *GroupMember) IsOwner() bool {
	return m.Role == MemberRoleOwner
}

func (m *GroupMember) IsAdmin() bool {
	return m.Role >= MemberRoleAdmin
}

func (m *GroupMember) SetRole(role MemberRole) {
	m.Role = role
}

type GroupJoinRequest struct {
	ID        string
	UserID    string
	GroupID   string
	Reason    string
	Status    RequestStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RequestStatus string

const (
	RequestStatusPending  RequestStatus = "pending"
	RequestStatusAccepted RequestStatus = "accepted"
	RequestStatusRejected RequestStatus = "rejected"
)

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

func (r *GroupJoinRequest) Accept() {
	r.Status = RequestStatusAccepted
	r.UpdatedAt = time.Now()
}

func (r *GroupJoinRequest) Reject() {
	r.Status = RequestStatusRejected
	r.UpdatedAt = time.Now()
}

func (r *GroupJoinRequest) IsPending() bool {
	return r.Status == RequestStatusPending
}
