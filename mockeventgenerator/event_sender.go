package main

import (
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

// Sender is an interface for sending events/messages
// The Send method receives a value of any type and returns an error

type Sender interface {
	Send(msg any) error
}

// RabbitMQSender implements Sender and sends messages to a RabbitMQ queue
// Usage: NewRabbitMQSender(url, queueName)
type RabbitMQSender struct {
	conn      *amqp091.Connection
	channel   *amqp091.Channel
	queueName string
}

func NewRabbitMQSender(url, queueName string) (*RabbitMQSender, error) {
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
	return &RabbitMQSender{conn: conn, channel: ch, queueName: queueName}, nil
}

func (s *RabbitMQSender) Send(msg any) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return s.channel.Publish(
		"", // exchange
		s.queueName,
		false, // mandatory
		false, // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (s *RabbitMQSender) Close() error {
	err1 := s.channel.Close()
	err2 := s.conn.Close()
	if err1 != nil {
		return err1
	}
	return err2
}
