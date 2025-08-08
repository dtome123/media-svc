package main

import (
	"log"
	"media-svc/config"
	"media-svc/internal/port"
	"media-svc/pkgs/rabbitmq"

	mongodb "github.com/dtome123/go-mongo-generic"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := mongodb.NewDatabase(
		mongodb.WithDatabase(cfg.DB.Mongo.Database),
		mongodb.WithSingleURL(cfg.DB.Mongo.DSN),
	)
	if err != nil {
		panic(err)
	}

	rabbitClient, err := rabbitmq.NewPublisher(cfg.RabbitMQ.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	server := port.NewServer(cfg, db, rabbitClient)
	server.Run()
}
