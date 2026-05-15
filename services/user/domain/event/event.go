package event

import "time"

type Event interface {
	EventType() string
	OccurredAt() time.Time
	AggregateID() string
}

type BaseEvent struct {
	eventType   string
	occurredAt  time.Time
	aggregateID string
}

func (e *BaseEvent) EventType() string {
	return e.eventType
}

func (e *BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e *BaseEvent) AggregateID() string {
	return e.aggregateID
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
	handlers, ok := p.handlers[event.EventType()]
	if !ok {
		return
	}
	for _, h := range handlers {
		go h(event)
	}
}
