package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"log-aggregator/internal/config"
	"log-aggregator/internal/model"

	"github.com/nats-io/nats.go"
)

func main() {
	cfg := config.LoadConfig()

	nc, err := nats.Connect(cfg.NATSUrl)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	log.Println("Producer started sending fake logs to NATS...")

	levels := []model.LogLevel{model.LevelInfo,model.LevelWarn,model.LevelError,}
	services := []string{"auth-service", "payment-service", "user-service"}

	for i := 1; i <= 50; i++ {
		logEntry := model.Log{
			Service:   services[rand.IntN(len(services))],
			Level:     levels[rand.IntN(len(levels))],
			Message:   fmt.Sprintf("Test log message number %d", i),
			Timestamp: time.Now(),
		}

		data, err := json.Marshal(logEntry)
		if err != nil {
			log.Printf("Error marshaling log: %v", err)
			continue
		}

		subject := fmt.Sprintf("logs.%s", logEntry.Service)
		if err := nc.Publish(subject, data); err != nil {
			log.Printf("Error publishing log: %v", err)
		} else {
			fmt.Printf("[%d/50] Sent log: %s -> %s\n", i, subject, logEntry.Message)
		}

		time.Sleep(50 * time.Millisecond) 
	}

	log.Println("All 50 test logs sent successfully!")
}