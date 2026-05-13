package entity

import "time"

type Group struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"owner_id"`
	Type        string    `json:"type"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GroupMember struct {
	GroupID  string    `json:"group_id"`
	UserID   string    `json:"user_id"`
	Role     int16     `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
	Nickname string    `json:"nickname"`
}

const (
	GroupRoleNormal = 1
	GroupRoleAdmin = 2
	GroupRoleOwner = 3
)

type GroupJoinRequest struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	GroupID   string    `json:"group_id"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	GroupJoinRequestStatusPending  = "pending"
	GroupJoinRequestStatusAccepted = "accepted"
	GroupJoinRequestStatusRejected = "rejected"
)