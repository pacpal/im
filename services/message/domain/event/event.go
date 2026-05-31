package event

import "IM/pkg/event"

type MessageSentEvent struct {
	event.BaseEvent
	MessageID  string
	SenderID   string
	ReceiverID string
	MsgType    string
}

type MessageReadEvent struct {
	event.BaseEvent
	MessageID string
	UserID    string
}

type MessageRevokedEvent struct {
	event.BaseEvent
	MessageID string
	UserID    string
}

type UserOnlineEvent struct {
	event.BaseEvent
	UserID string
}

type UserOfflineEvent struct {
	event.BaseEvent
	UserID string
}