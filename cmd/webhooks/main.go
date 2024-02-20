package main

import (
	"context"
	"os"

	"github.com/gavsidhu/notion-hooks/internal/config"
	"github.com/gavsidhu/notion-hooks/internal/logging"
	"github.com/gavsidhu/notion-hooks/internal/webhook"
	"github.com/gavsidhu/notion-hooks/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var maxWorkers = 1

func main() {
	logging.Logger.Info("Starting the application")
	err := godotenv.Load()
	if err != nil {
		logging.Logger.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	rabbitMQ, err := config.NewRabbitMQConnection(os.Getenv("RABBITMQ_CONNECTION_URL"))
	if err != nil {
		logging.Logger.Fatal(err)
	}

	defer rabbitMQ.Close()

	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		logging.Logger.Fatal(err)
	}

	go webhook.StartPollingDatabase(ctx, dbpool, rabbitMQ.Ch)

	for i := 0; i < maxWorkers; i++ {
		go worker.StartWorker(rabbitMQ, "proccessingQueue", dbpool, webhook.ProccessWebhook)

	}

	for i := 0; i < maxWorkers; i++ {
		go worker.StartWorker(rabbitMQ, "eventsQueue", dbpool, webhook.SendEventsToUser)
	}

	for i := 0; i < maxWorkers; i++ {
		go worker.StartWorker(rabbitMQ, "initalPollQueue", dbpool, webhook.HandleInitialPolling)
	}

	select {}

}
