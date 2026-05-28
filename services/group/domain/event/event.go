// Package event 定义群组域的领域事件及简单的事件发布器实现。
package event

import "time"

// Event 是领域事件的接口。
type Event interface {
	GetEventType() string
	GetOccurredAt() time.Time
	GetAggregateID() string
}

// BaseEvent 提供事件的基础字段与默认实现。
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

// 具体事件类型（创建、更新、成员变更等）。
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
	GroupID string
	UserID  string
	AdminID string
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
	GroupID    string
	OldOwnerID string
	NewOwnerID string
}

// Handler 是事件处理函数签名。
type Handler func(event Event)

// Publisher 定义事件发布器接口。
type Publisher interface {
	Subscribe(eventType string, handler Handler)
	Publish(event Event)
}

// EventPublisher 是一个简单的内存事件总线实现。
type EventPublisher struct {
	handlers map[string][]Handler
}

// NewEventPublisher 创建一个 EventPublisher。
func NewEventPublisher() *EventPublisher {
	return &EventPublisher{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe 注册事件处理器。
func (p *EventPublisher) Subscribe(eventType string, handler Handler) {
	p.handlers[eventType] = append(p.handlers[eventType], handler)
}

// Publish 发布事件，会异步执行已注册的处理器。
func (p *EventPublisher) Publish(event Event) {
	handlers, ok := p.handlers[event.GetEventType()]
	if !ok {
		return
	}
	for _, h := range handlers {
		go h(event)
	}
}
