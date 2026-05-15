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

type UserRegisteredEvent struct {
	BaseEvent
	UserID   string
	Name     string
	Tele     string
}

type UserLoggedInEvent struct {
	BaseEvent
	UserID string
	Tele   string
}

type FriendRequestCreatedEvent struct {
	BaseEvent
	RequestID string
	FromUID   string
	ToUID     string
}

type FriendRequestAcceptedEvent struct {
	BaseEvent
	RequestID string
	FromUID   string
	ToUID     string
}

type FriendRequestRejectedEvent struct {
	BaseEvent
	RequestID string
	FromUID   string
	ToUID     string
}

type FriendshipCreatedEvent struct {
	BaseEvent
	UserID   string
	FriendID string
}

type FriendshipDeletedEvent struct {
	BaseEvent
	UserID   string
	FriendID string
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
