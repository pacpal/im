// Package event 定义消息服务的领域事件与简单事件发布器实现。
package event

import "time"

// Event 是领域事件接口。
type Event interface {
	GetEventType() string
	GetOccurredAt() time.Time
	GetAggregateID() string
}

// BaseEvent 提供事件通用字段与默认实现。
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

// 具体事件类型。
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

// Handler 事件处理函数类型。
type Handler func(event Event)

// Publisher 定义事件发布器接口。
type Publisher interface {
	Subscribe(eventType string, handler Handler)
	Publish(event Event)
}

// EventPublisher 是内存级别的简单事件总线，实现 Publish/Subscribe。
type EventPublisher struct {
	handlers map[string][]Handler
}

// NewEventPublisher 创建 EventPublisher 实例。
func NewEventPublisher() *EventPublisher {
	return &EventPublisher{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe 注册指定事件类型的处理器。
func (p *EventPublisher) Subscribe(eventType string, handler Handler) {
	p.handlers[eventType] = append(p.handlers[eventType], handler)
}

// Publish 发布事件，已注册的处理器将被异步调用。
func (p *EventPublisher) Publish(event Event) {
	handlers, ok := p.handlers[event.GetEventType()]
	if !ok {
		return
	}
	for _, h := range handlers {
		go h(event)
	}
}
