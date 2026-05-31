package event

import (
	"IM/pkg/logger"
	"IM/pkg/mq"
	"encoding/json"
	"fmt"
)

type EventRouter struct {
	handlers map[string][]func(data []byte) error
}

func NewEventRouter() *EventRouter {
	return &EventRouter{
		handlers: make(map[string][]func(data []byte) error),
	}
}

func (r *EventRouter) Register(eventType string, handler func(data []byte) error) {
	r.handlers[eventType] = append(r.handlers[eventType], handler)
}

func (r *EventRouter) Handle(eventType string, data []byte) error {
	handlers, ok := r.handlers[eventType]
	if !ok {
		logger.Debugw("No handler registered for event", "component", "event_router", "event_type", eventType)
		return nil
	}

	var lastErr error
	for _, handler := range handlers {
		if err := handler(data); err != nil {
			logger.Errorw("Event handler failed", "component", "event_router", "event_type", eventType, "err", err)
			lastErr = err
		}
	}
	return lastErr
}

type RabbitMQPublisher struct {
	producer *mq.Producer
	exchange string
}

func NewRabbitMQPublisher(conn *mq.Connection, exchange string) (*RabbitMQPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	if err := mq.DeclareExchange(ch, exchange, "topic"); err != nil {
		ch.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}
	ch.Close()

	producer, err := mq.NewProducer(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &RabbitMQPublisher{
		producer: producer,
		exchange: exchange,
	}, nil
}

func (p *RabbitMQPublisher) Publish(event Event) {
	body, err := json.Marshal(event)
	if err != nil {
		logger.Errorw("Failed to marshal event", "component", "event_publisher", "err", err)
		return
	}

	routingKey := event.GetEventType()
	if err := p.producer.Publish(p.exchange, routingKey, body); err != nil {
		logger.Errorw("Failed to publish event", "component", "event_publisher", "err", err, "event_type", routingKey)
		return
	}

	logger.Infow("Event published", "component", "event_publisher", "event_type", routingKey, "aggregate_id", event.GetAggregateID())
}

func (p *RabbitMQPublisher) Close() error {
	return p.producer.Close()
}

type EventSubscriber struct {
	consumer  *mq.Consumer
	queueName string
	router    *EventRouter
}

func NewEventSubscriber(conn *mq.Connection, exchange, queueName string, routingKeys []string) (*EventSubscriber, error) {
	consumer, err := mq.NewConsumer(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	if err := consumer.DeclareQueue(queueName, true); err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	for _, routingKey := range routingKeys {
		if err := consumer.BindQueue(queueName, routingKey, exchange); err != nil {
			return nil, fmt.Errorf("failed to bind queue for %s: %w", routingKey, err)
		}
	}

	router := NewEventRouter()

	return &EventSubscriber{
		consumer:  consumer,
		queueName: queueName,
		router:    router,
	}, nil
}

func (s *EventSubscriber) Register(eventType string, handler func(data []byte) error) {
	s.router.Register(eventType, handler)
}

func (s *EventSubscriber) Start() error {
	return s.consumer.Consume(s.queueName, func(body []byte) error {
		var baseEvent BaseEvent
		if err := json.Unmarshal(body, &baseEvent); err != nil {
			logger.Errorw("Failed to unmarshal event", "component", "event_subscriber", "err", err)
			return err
		}

		eventType := baseEvent.GetEventType()
		logger.Infow("Event received", "component", "event_subscriber", "event_type", eventType)

		return s.router.Handle(eventType, body)
	})
}

func (s *EventSubscriber) Close() error {
	return s.consumer.Close()
}

type EventBus struct {
	localPublisher    *LocalEventBus
	rabbitmqPublisher *RabbitMQPublisher
}

func NewEventBus(localPublisher *LocalEventBus, rabbitmqPublisher *RabbitMQPublisher) *EventBus {
	return &EventBus{
		localPublisher:    localPublisher,
		rabbitmqPublisher: rabbitmqPublisher,
	}
}

func (bus *EventBus) Publish(event Event) {
	bus.localPublisher.Publish(event)

	if bus.rabbitmqPublisher != nil {
		bus.rabbitmqPublisher.Publish(event)
	}
}

func (bus *EventBus) Subscribe(eventType string, handler Handler) {
	bus.localPublisher.Subscribe(eventType, handler)
}
