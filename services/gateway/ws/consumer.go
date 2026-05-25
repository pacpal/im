// Package ws 提供消息消费者，将后端 MQ 消息推送给在线 WebSocket 客户端。
package ws

import (
	"IM/pkg/logger"
	"context"
	"encoding/json"
	"time"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// PushMessage 是网关推送给客户端的消息结构。
type PushMessage struct {
	MsgID      string `json:"msg_id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
	MsgType    string `json:"msg_type"`
	Timestamp  int64  `json:"timestamp"`
	IsGroup    bool   `json:"is_group"`
}

// MessageConsumer 从 RabbitMQ 消费消息并推送给 WebSocket 用户。
type MessageConsumer struct {
	hub       *Hub
	conn      *amqp091.Connection
	channel   *amqp091.Channel
	queueName string
	exchange  string
}

// NewMessageConsumer 创建并返回 MessageConsumer。
func NewMessageConsumer(hub *Hub, amqpURL, exchange, queueName string) (*MessageConsumer, error) {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	if _, err := ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	// 绑定通配符路由键，接收所有用户消息
	if err := ch.QueueBind(queueName, "message.*", exchange, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &MessageConsumer{
		hub:       hub,
		conn:      conn,
		channel:   ch,
		queueName: queueName,
		exchange:  exchange,
	}, nil
}

// Start 启动后台协程消费消息。
func (c *MessageConsumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(c.queueName, "", false, false, false, false, nil)
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
				c.handleMessage(msg)
			}
		}
	}()

	return nil
}

func (c *MessageConsumer) handleMessage(msg amqp091.Delivery) {
	var pushMsg PushMessage
	if err := json.Unmarshal(msg.Body, &pushMsg); err != nil {
		logger.Warnw("Failed to unmarshal push message", "component", "gateway_consumer", "err", err)
		msg.Nack(false, false)
		return
	}

	data, err := json.Marshal(pushMsg)
	if err != nil {
		msg.Nack(false, false)
		return
	}

	// 如果用户在线则推送
	if c.hub.IsOnline(pushMsg.ReceiverID) {
		c.hub.SendToUser(pushMsg.ReceiverID, data)
		logger.Infow("Push message to user", "component", "gateway_consumer", "user", pushMsg.ReceiverID, "msg_id", pushMsg.MsgID)
		msg.Ack(false)
	} else {
		// 用户不在线，不处理，让消息服务通过离线消息处理
		msg.Nack(false, true)
	}
}

// Close 关闭消费者连接。
func (c *MessageConsumer) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// WaitForRabbitMQ 简单重试等待 RabbitMQ 可用。
func WaitForRabbitMQ(amqpURL string, maxRetries int, interval time.Duration) (*amqp091.Connection, error) {
	var conn *amqp091.Connection
	var err error
	for i := 0; i < maxRetries; i++ {
		conn, err = amqp091.Dial(amqpURL)
		if err == nil {
			return conn, nil
		}
		time.Sleep(interval)
	}
	return nil, err
}
