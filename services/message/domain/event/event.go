package event

import "time"

type Event interface {
	GetEventType() string
	GetOccurredAt() time.Time
	GetAggregateID() string
}

type BaseEvent struct {
	EventType   string
	OccurredAt  time.Time
	AggregateID string
}

func (e *BaseEvent) GetEventType() string {
	return e.EventType
}

func (e *BaseEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

func (e *BaseEvent) GetAggregateID() string {
	return e.AggregateID
}

type MessageSentEvent struct {
	BaseEvent
	MessageID  string
	SenderID   string
	ReceiverID string
	MsgType    string
}

type MessageReadEvent struct {
	BaseEvent
	MessageID string
	UserID    string
}

type MessageRevokedEvent struct {
	BaseEvent
	MessageID string
	UserID    string
}

type UserOnlineEvent struct {
	BaseEvent
	UserID string
}

type UserOfflineEvent struct {
	BaseEvent
	UserID string
}

type Handler func(event Event)

type Publisher interface {
	Subscribe(eventType string, handler Handler)
	Publish(event Event)
}

type EventPublisher struct {
	handlers map[string][]Handler
}

func NewEventPublisher() *EventPublisher {
	return &EventPublisher{
		handlers: make(map[string][]Handler),
	}
}

func (p *EventPublisher) Subscribe(eventType string, handler Handler) {
	p.handlers[eventType] = append(p.handlers[eventType], handler)
}

func (p *EventPublisher) Publish(event Event) {
	handlers, ok := p.handlers[event.GetEventType()]
	if !ok {
		return
	}
	for _, h := range handlers {
		go h(event)
	}
}
