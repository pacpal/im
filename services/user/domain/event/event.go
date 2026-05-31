package event

import "IM/pkg/event"

type UserRegisteredEvent struct {
	event.BaseEvent
	UserID string
	Name   string
	Tele   string
}

type UserLoggedInEvent struct {
	event.BaseEvent
	UserID string
	Tele   string
}

type FriendRequestCreatedEvent struct {
	event.BaseEvent
	RequestID string
	FromUID   string
	ToUID     string
}

type FriendRequestAcceptedEvent struct {
	event.BaseEvent
	RequestID string
	FromUID   string
	ToUID     string
}

type FriendRequestRejectedEvent struct {
	event.BaseEvent
	RequestID string
	FromUID   string
	ToUID     string
}

type FriendshipCreatedEvent struct {
	event.BaseEvent
	UserID   string
	FriendID string
}

type FriendshipDeletedEvent struct {
	event.BaseEvent
	UserID   string
	FriendID string
}