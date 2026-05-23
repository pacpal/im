// Package event 定义领域事件及简单的发布/订阅机制，用于服务内事件分发。
package event

import "time"

// Event 是领域事件的通用接口。
type Event interface {
	GetEventType() string
	GetOccurredAt() time.Time
	GetAggregateID() string
}

// BaseEvent 提供事件公共字段的基础实现。
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

// 具体事件类型示例（注册、登录、好友请求等）。
type UserRegisteredEvent struct {
	BaseEvent
	UserID string
	Name   string
	Tele   string
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

// Handler 是事件处理函数签名。
type Handler func(event Event)

// Publisher 定义事件发布器接口。
type Publisher interface {
	Subscribe(eventType string, handler Handler)
	Publish(event Event)
}

// EventPublisher 是一个轻量的内存事件总线实现。
type EventPublisher struct {
	handlers map[string][]Handler
}

// NewEventPublisher 创建一个新的事件发布器。
func NewEventPublisher() *EventPublisher {
	return &EventPublisher{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe 注册事件处理函数到指定事件类型。
func (p *EventPublisher) Subscribe(eventType string, handler Handler) {
	p.handlers[eventType] = append(p.handlers[eventType], handler)
}

// Publish 发布事件，会异步调用已注册的处理函数。
func (p *EventPublisher) Publish(event Event) {
	handlers, ok := p.handlers[event.GetEventType()]
	if !ok {
		return
	}
	for _, h := range handlers {
		go h(event)
	}
}
