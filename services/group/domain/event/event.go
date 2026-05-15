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

type GroupCreatedEvent struct {
	BaseEvent
	GroupID string
	Name    string
	OwnerID string
}

type GroupUpdatedEvent struct {
	BaseEvent
	GroupID string
	Name    string
}

type GroupDeletedEvent struct {
	BaseEvent
	GroupID string
	OwnerID string
}

type MemberJoinedEvent struct {
	BaseEvent
	GroupID string
	UserID  string
	Role    int
}

type MemberLeftEvent struct {
	BaseEvent
	GroupID string
	UserID  string
}

type MemberRemovedEvent struct {
	BaseEvent
	GroupID  string
	UserID   string
	OwnerID  string
}

type JoinRequestCreatedEvent struct {
	BaseEvent
	RequestID string
	UserID    string
	GroupID   string
}

type JoinRequestAcceptedEvent struct {
	BaseEvent
	RequestID string
	UserID    string
	GroupID   string
}

type JoinRequestRejectedEvent struct {
	BaseEvent
	RequestID string
	UserID    string
	GroupID   string
}

type OwnerTransferredEvent struct {
	BaseEvent
	GroupID     string
	OldOwnerID  string
	NewOwnerID  string
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
