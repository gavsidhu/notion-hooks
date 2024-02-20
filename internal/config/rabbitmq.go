package config

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQConnection struct {
	Conn *amqp091.Connection
	Ch   *amqp091.Channel
}

func NewRabbitMQConnection(url string) (*RabbitMQConnection, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Error in channel creation")
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare("proccessingQueue", true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare("eventsQueue", true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare("initalPollQueue", true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQConnection{
		Conn: conn,
		Ch:   ch,
	}, nil
}

func (rmq *RabbitMQConnection) Close() {
	if rmq.Ch != nil {
		rmq.Ch.Close()
	}
	if rmq.Conn != nil {
		rmq.Conn.Close()
	}
}
