package mq

import (
	"IM/services/message/domain/entity"
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConnection struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

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

func (r *RabbitMQConnection) Close() error {
	if r.ch != nil {
		r.ch.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

func (r *RabbitMQConnection) GetChannel() *amqp091.Channel {
	return r.ch
}

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

func BindQueue(ch *amqp091.Channel, queueName, routingKey, exchange string) error {
	return ch.QueueBind(
		queueName,
		routingKey,
		exchange,
		false,
		nil,
	)
}
