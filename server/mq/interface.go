// Package mq maybe RabbitMQ
package mq

type Producer interface {
	Publish(topic string, payload []byte) error
}
