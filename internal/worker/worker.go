package worker

import (
	"fmt"
	"log"

	"github.com/gavsidhu/notion-hooks/internal/config"
	"github.com/gavsidhu/notion-hooks/internal/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
)

type handlerFunc func(msg amqp091.Delivery, ch *amqp091.Channel, pool *pgxpool.Pool)

func StartWorker(rabbitMQ *config.RabbitMQConnection, queueName string, pool *pgxpool.Pool, handler handlerFunc) {
	logging.Logger.Info(fmt.Sprintf("Starting worker for queue: %s", queueName))

	ch, err := rabbitMQ.Conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		logging.Logger.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		handler(msg, ch, pool)
	}
}
