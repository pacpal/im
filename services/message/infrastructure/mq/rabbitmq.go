// Package mq 提供 message 服务使用的 RabbitMQ 封装，支持自动重连。
package mq

import (
	"IM/pkg/logger"
	"IM/services/message/domain/entity"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// RabbitMQConnection 包装底层连接和通道，支持自动重连。
type RabbitMQConnection struct {
	url      string
	conn     *amqp091.Connection
	ch       *amqp091.Channel
	mu       sync.RWMutex
	closed   bool
	closeCh  chan struct{}
	notifyCh <-chan *amqp091.Error
}

// NewRabbitMQConnection 连接到 RabbitMQ 并返回包装对象。
func NewRabbitMQConnection(url string) (*RabbitMQConnection, error) {
	r := &RabbitMQConnection{
		url:     url,
		closeCh: make(chan struct{}),
	}
	if err := r.connect(); err != nil {
		return nil, err
	}
	go r.handleReconnect()
	return r, nil
}

func (r *RabbitMQConnection) connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	conn, err := amqp091.Dial(r.url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	r.conn = conn
	r.ch = ch
	r.notifyCh = conn.NotifyClose(make(chan *amqp091.Error, 1))
	return nil
}

// handleReconnect 监听连接关闭事件，断开时自动重连。
func (r *RabbitMQConnection) handleReconnect() {
	for {
		select {
		case <-r.closeCh:
			return
		case err := <-r.notifyCh:
			if err == nil {
				return
			}
			logger.Warnw("RabbitMQ connection closed, reconnecting...", "component", "mq", "err", err)

			r.mu.RLock()
			closed := r.closed
			r.mu.RUnlock()
			if closed {
				return
			}

			for i := 0; i < 10; i++ {
				time.Sleep(time.Duration(i+1) * time.Second)
				if err := r.connect(); err != nil {
					logger.Warnw("RabbitMQ reconnect failed", "component", "mq", "err", err, "retry", i+1)
					continue
				}
				logger.Infow("RabbitMQ reconnected successfully", "component", "mq")
				break
			}
		}
	}
}

// Close 关闭通道与连接。
func (r *RabbitMQConnection) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.closed = true
	close(r.closeCh)

	var errs []error
	if r.ch != nil {
		if err := r.ch.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// GetChannel 返回底层 amqp Channel（线程安全）。
func (r *RabbitMQConnection) GetChannel() *amqp091.Channel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ch
}

// MessageProducer 将消息序列化并发布到指定交换器。
type MessageProducer struct {
	conn     *RabbitMQConnection
	exchange string
}

func NewMessageProducer(conn *RabbitMQConnection, exchange string) *MessageProducer {
	return &MessageProducer{
		conn:     conn,
		exchange: exchange,
	}
}

// PublishMessage 发布消息到以 receiverID 为路由键的 topic。
func (p *MessageProducer) PublishMessage(ctx context.Context, msg *entity.Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ch := p.conn.GetChannel()
	if ch == nil {
		return fmt.Errorf("channel is not available")
	}

	routingKey := "message." + msg.ReceiverID
	err = ch.PublishWithContext(ctx,
		p.exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:        body,
			DeliveryMode: amqp091.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

// MessageConsumer 从队列消费并调用 handler 处理消息。
type MessageConsumer struct {
	conn      *RabbitMQConnection
	queueName string
}

func NewMessageConsumer(conn *RabbitMQConnection, queueName string) *MessageConsumer {
	return &MessageConsumer{
		conn:      conn,
		queueName: queueName,
	}
}

// Consume 开始消费队列并在后台协程调用 handler 处理消息，失败时根据返回值决定是否重排。
func (c *MessageConsumer) Consume(ctx context.Context, handler func(*entity.Message) error) error {
	ch := c.conn.GetChannel()
	if ch == nil {
		return fmt.Errorf("channel is not available")
	}

	msgs, err := ch.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				var message entity.Message
				if err := json.Unmarshal(msg.Body, &message); err != nil {
					msg.Nack(false, false)
					continue
				}

				if err := handler(&message); err != nil {
					msg.Nack(false, true)
				} else {
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

// DeclareExchange 声明交换器（topic）。
func DeclareExchange(ch *amqp091.Channel, name string) error {
	return ch.ExchangeDeclare(
		name,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

// DeclareQueue 声明持久化队列。
func DeclareQueue(ch *amqp091.Channel, name string) error {
	_, err := ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

// BindQueue 绑定队列到交换器与路由键。
func BindQueue(ch *amqp091.Channel, queueName, routingKey, exchange string) error {
	return ch.QueueBind(
		queueName,
		routingKey,
		exchange,
		false,
		nil,
	)
}
