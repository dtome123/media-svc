package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"media-svc/config"
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
	client, err := rabbitmq.New(cfg.RabbitMQ.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Initialize service layer
	service := services.NewService(cfg, db, client)

	// Create context with cancel to control lifecycle
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Start consumer goroutine to process jobs from queue
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := client.Consume(ctx, cfg.RabbitMQ.Queue, func(data []byte) error {
			var job types.TranscodeJob

			// Parse job from JSON payload
			if err := json.Unmarshal(data, &job); err != nil {
				return fmt.Errorf("invalid job format: %w", err)
			}

			log.Printf("Received job: %+v", job)

			// Call transcoding service
			_, err := service.GetMediaSvc().TranscodeVideo(context.Background(), media.TranscodeVideoInput{
				FilePath: job.InputPath,
			})
			if err != nil {
				return fmt.Errorf("transcode failed: %w", err)
			}

			return nil
		})

		// Log when consumer exits
		if err != nil {
			log.Printf("Consumer exited with error: %v", err)
		} else {
			log.Println("Consumer exited cleanly")
		}
	}()

	// Listen for OS signals for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal received
	<-sig
	log.Println("Worker shutting down...")

	// Cancel context to stop consumer
	cancel()

	// Wait for consumer goroutine to finish
	wg.Wait()

	// Close RabbitMQ connection
	client.Close()

	log.Println("Worker stopped.")
}
