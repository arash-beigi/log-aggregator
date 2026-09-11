package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"log-aggregator/internal/config"
	"log-aggregator/internal/model"
	"log-aggregator/internal/worker"

	"github.com/nats-io/nats.go"
)

type Consumer struct {
	nc   *nats.Conn
	js   nats.JetStreamContext
	sub  *nats.Subscription
	pool *worker.Pool
	cfg  *config.Config
}

func NewConsumer(cfg *config.Config, pool *worker.Pool) (*Consumer, error) {
	nc, err := nats.Connect(cfg.NATSUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create jetstream context: %w", err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "LOGS",
		Subjects: []string{"logs.>"},
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		nc.Close()
		return nil, fmt.Errorf("failed to create nats stream: %w", err)
	}

	return &Consumer{
		nc:   nc,
		js:   js,
		pool: pool,
		cfg:  cfg,
	}, nil
}

func (p *Consumer) Start(ctx context.Context) error {
	sub, err := p.js.Subscribe("logs.>", func(msg *nats.Msg) {
		var logItem model.Log

		if err := json.Unmarshal(msg.Data, &logItem); err != nil {
			log.Printf("[NATS Consumer] Error unmarshaling log message: %v", err)
			_ = msg.Nak()
			return
		}

		p.pool.Submit(logItem)

		if err := msg.Ack(); err != nil {
			log.Printf("[NATS Consumer] Error sending Ack: %v", err)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to nats subject: %w", err)
	}

	p.sub = sub
	log.Println("NATS Consumer successfully subscribed to subject 'logs.>'")
	return nil
}

func (p *Consumer) Close() {
	if p.sub != nil {
		_ = p.sub.Unsubscribe()
	}
	if p.nc != nil {
		p.nc.Close()
	}
	log.Println("NATS Consumer connection closed")
}
