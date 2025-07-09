package main

import (
	"github.com/rabbitmq/amqp091-go"
)

type Receiver interface {
	Receive(handler func(body []byte, ack func()) error) error
}

// RabbitMQReceiver implements Receiver and receives messages from a RabbitMQ queue
// Usage: NewRabbitMQReceiver(url, queueName)
type RabbitMQReceiver struct {
	conn      *amqp091.Connection
	channel   *amqp091.Channel
	queueName string
}

func NewRabbitMQReceiver(url, queueName string) (*RabbitMQReceiver, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}
	return &RabbitMQReceiver{conn: conn, channel: ch, queueName: queueName}, nil
}

// Receive consumes messages from the queue and calls handler for each message body
func (r *RabbitMQReceiver) Receive(handler func(body []byte, ackEventFunc func()) error) error {
	deliveries, err := r.channel.Consume(
		r.queueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}
	for msg := range deliveries {
		ackEventFunc := func() { msg.Ack(false) }
		if err := handler(msg.Body, ackEventFunc); err != nil {
			return err
		}
	}
	return nil
}

func (r *RabbitMQReceiver) Close() error {
	err1 := r.channel.Close()
	err2 := r.conn.Close()
	if err1 != nil {
		return err1
	}
	return err2
}
