package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log-aggregator/internal/config"
	"log-aggregator/internal/nats"
	"log-aggregator/internal/repository"
	"log-aggregator/internal/worker"
)

func main() {
	log.Println("Starting Log Aggregator Consumer Service...")

	cfg := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	
	mongoRepo, err := repository.NewMongoRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Mongo Repository: %v", err)
	}
	defer func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := mongoRepo.Close(shutdownCtx); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}()

	
	bufferSize := 5000
	pool := worker.NewPool(mongoRepo, cfg, bufferSize)
	pool.Start(ctx)


	consumer, err := nats.NewConsumer(cfg, pool)
	if err != nil {
		log.Fatalf("Failed to initialize NATS Consumer: %v", err)
	}
	defer consumer.Close()

	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("Failed to start NATS Consumer: %v", err)
	}

	log.Println("Log Aggregator Consumer is running and waiting for logs...")

	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, closing connections...")

	consumer.Close()
	pool.Stop()

	log.Println("Log Aggregator Consumer stopped cleanly.")
}