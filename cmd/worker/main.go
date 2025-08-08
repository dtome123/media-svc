package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"media-svc/config"
	"media-svc/internal/job/transcode"
	"media-svc/internal/services"
	"media-svc/internal/services/media"
	"media-svc/internal/types"
	"media-svc/pkgs/rabbitmq"
	"os"
	"os/signal"
	"sync"
	"syscall"

	mongodb "github.com/dtome123/go-mongo-generic"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Initialize MongoDB connection
	db, err := mongodb.NewDatabase(
		mongodb.WithDatabase(cfg.DB.Mongo.Database),
		mongodb.WithSingleURL(cfg.DB.Mongo.DSN),
	)
	if err != nil {
		panic(err)
	}

	// Initialize RabbitMQ client
	client, err := rabbitmq.NewConsumer(cfg.RabbitMQ.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Initialize service layer
	service := services.NewService(cfg, db, nil)

	// Create orchestrator with desired number of workers and queue depth
	workerCount := 4
	queueDepth := 100
	orcTranscode := transcode.New(service, workerCount, queueDepth)

	// Start the orchestrator (start worker goroutines)
	orcTranscode.Start()

	// Context & WaitGroup to manage lifecycle of RabbitMQ consumer
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Start consumer goroutine to consume from RabbitMQ
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := client.Consume(ctx, cfg.RabbitMQ.Queue, func(data []byte) error {
			var job types.TranscodeJob
			if err := json.Unmarshal(data, &job); err != nil {
				return fmt.Errorf("invalid job format: %w", err)
			}

			log.Printf("Received job: %+v", job)
			orcTranscode.AddJob(media.TranscodeVideoInput{MediaID: job.MediaID})

			return nil
		})

		if err != nil {
			// Nếu lỗi là channel closed do ctx bị cancel thì có thể không cần log lỗi nặng
			if err.Error() == "channel closed" {
				log.Println("Consumer stopped due to channel close (likely normal shutdown)")
			} else {
				log.Printf("Consumer exited with error: %v", err)
			}
		} else {
			log.Println("Consumer exited cleanly")
		}
	}()

	// Listen for OS signals for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	log.Println("Worker shutting down...")

	// Stop consumer
	cancel()
	wg.Wait()

	// Stop orchestrator and wait for all workers to finish
	orcTranscode.Stop()

	// Close RabbitMQ client connection
	client.Close()

	log.Println("Worker stopped.")
}
