// Package mq 提供 message 服务使用的 RabbitMQ 封装。
package mq

import (
	"IM/services/message/domain/entity"
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// RabbitMQConnection 包装底层连接和通道，便于管理。
type RabbitMQConnection struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

// NewRabbitMQConnection 连接到 RabbitMQ 并返回包装对象。
func NewRabbitMQConnection(url string) (*RabbitMQConnection, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQConnection{conn: conn, ch: ch}, nil
}

// Close 关闭通道与连接。
func (r *RabbitMQConnection) Close() error {
	if r.ch != nil {
		r.ch.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// GetChannel 返回底层 amqp Channel。
func (r *RabbitMQConnection) GetChannel() *amqp091.Channel {
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

	routingKey := "message." + msg.ReceiverID
	return p.conn.ch.PublishWithContext(ctx,
		p.exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
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
	msgs, err := c.conn.ch.Consume(
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
