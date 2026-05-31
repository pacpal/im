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

func (e *BaseEvent) GetEventType() string     { return e.EventType }
func (e *BaseEvent) GetOccurredAt() time.Time { return e.OccurredAt }
func (e *BaseEvent) GetAggregateID() string   { return e.AggregateID }

type Handler func(event Event)

type Publisher interface {
	Subscribe(eventType string, handler Handler)
	Publish(event Event)
}
