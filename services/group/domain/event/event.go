package event

import "IM/pkg/event"

type GroupCreatedEvent struct {
	event.BaseEvent
	GroupID string
	Name    string
	OwnerID string
}

type GroupUpdatedEvent struct {
	event.BaseEvent
	GroupID string
	Name    string
}

type GroupDeletedEvent struct {
	event.BaseEvent
	GroupID string
	OwnerID string
}

type MemberJoinedEvent struct {
	event.BaseEvent
	GroupID string
	UserID  string
	Role    int
}

type MemberLeftEvent struct {
	event.BaseEvent
	GroupID string
	UserID  string
}

type ChangedEvent struct {
	event.BaseEvent
	GroupID string
	UserID  string
	AdminID string
}

type MemberRemovedEvent struct {
	event.BaseEvent
	GroupID string
	UserID  string
	AdminID string
}

type JoinRequestCreatedEvent struct {
	event.BaseEvent
	RequestID string
	UserID    string
	GroupID   string
}

type JoinRequestAcceptedEvent struct {
	event.BaseEvent
	RequestID string
	UserID    string
	GroupID   string
}

type JoinRequestRejectedEvent struct {
	event.BaseEvent
	RequestID string
	UserID    string
	GroupID   string
}

type OwnerTransferredEvent struct {
	event.BaseEvent
	GroupID    string
	OldOwnerID string
	NewOwnerID string
}