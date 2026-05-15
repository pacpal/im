package mq

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// Connection RabbitMQ 连接管理
type Connection struct {
	conn *amqp091.Connection
	url  string
	mu   sync.RWMutex
}

// NewConnection 创建 RabbitMQ 连接
func NewConnection(url string) (*Connection, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return &Connection{
		conn: conn,
		url:  url,
	}, nil
}

// Close 关闭连接
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil && !c.conn.IsClosed() {
		return c.conn.Close()
	}
	return nil
}

// Channel 获取新 channel
func (c *Connection) Channel() (*amqp091.Channel, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.conn == nil || c.conn.IsClosed() {
		return nil, fmt.Errorf("connection is closed")
	}
	return c.conn.Channel()
}

// Producer 消息生产者
type Producer struct {
	conn *Connection
	ch   *amqp091.Channel
}

// NewProducer 创建消息生产者
func NewProducer(conn *Connection) (*Producer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}
	return &Producer{
		conn: conn,
		ch:   ch,
	}, nil
}

// Publish 发布消息到 exchange
func (p *Producer) Publish(exchange, routingKey string, body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return p.ch.PublishWithContext(ctx,
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp091.Publishing{
			ContentType: "application/octet-stream",
			Body:        body,
		},
	)
}

// Close 关闭生产者 channel
func (p *Producer) Close() error {
	if p.ch != nil {
		return p.ch.Close()
	}
	return nil
}

// Consumer 消息消费者
type Consumer struct {
	conn *Connection
	ch   *amqp091.Channel
}

// NewConsumer 创建消息消费者
func NewConsumer(conn *Connection) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}
	return &Consumer{
		conn: conn,
		ch:   ch,
	}, nil
}

// DeclareQueue 声明队列
func (c *Consumer) DeclareQueue(queueName string, durable bool) error {
	_, err := c.ch.QueueDeclare(
		queueName,
		durable,
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	return err
}

// BindQueue 绑定队列到 exchange
func (c *Consumer) BindQueue(queueName, routingKey, exchange string) error {
	return c.ch.QueueBind(
		queueName,
		routingKey,
		exchange,
		false, // noWait
		nil,   // args
	)
}

// Consume 消费消息
func (c *Consumer) Consume(queueName string, handler func([]byte) error) error {
	msgs, err := c.ch.Consume(
		queueName,
		"",    // consumer tag
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				_ = msg.Nack(false, true) // requeue on failure
			} else {
				_ = msg.Ack(false)
			}
		}
	}()

	return nil
}

// Close 关闭消费者 channel
func (c *Consumer) Close() error {
	if c.ch != nil {
		return c.ch.Close()
	}
	return nil
}

// DeclareExchange 声明交换机
func DeclareExchange(ch *amqp091.Channel, name, kind string) error {
	return ch.ExchangeDeclare(
		name,
		kind,
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,   // args
	)
}
