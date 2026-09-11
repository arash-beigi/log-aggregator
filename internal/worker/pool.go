package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"log-aggregator/internal/config"
	"log-aggregator/internal/model"
	"log-aggregator/internal/repository"
)

type Pool struct {
	repo        repository.LogRepository
	logChan     chan model.Log
	workerCount int
	batchSize   int
	flushPeriod time.Duration
	wg          sync.WaitGroup
}

func NewPool(repo repository.LogRepository, cfg *config.Config, bufferSize int) *Pool {
	return &Pool{
		repo:        repo,
		logChan:     make(chan model.Log, bufferSize),
		workerCount: cfg.WorkerCount,
		batchSize:   cfg.BatchSize,
		flushPeriod: 3 * time.Second,
	}
}
func (p *Pool) Start(ctx context.Context) {
	for i := 1; i <= p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}
	log.Printf("Worker Pool started with %d workers", p.workerCount)
}

func (p *Pool) Submit(logItem model.Log) {
	p.logChan <- logItem
}

func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()

	batch := make([]model.Log, 0, p.batchSize)

	ticker := time.NewTicker(p.flushPeriod)
	defer ticker.Stop()

	for {
		select {
		case logItem, ok := <-p.logChan:
			if !ok {
				p.flush(ctx, id, &batch)
				return
			}
			batch = append(batch, logItem)

			if len(batch) >= p.batchSize {
				p.flush(ctx, id, &batch)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				p.flush(ctx, id, &batch)
			}
		case <-ctx.Done():
			p.flush(ctx, id, &batch)
			return
		}
	}
}

func (p *Pool) flush(ctx context.Context, workerID int, batch *[]model.Log) {
	if len(*batch) == 0 {
		return
	}

	err := p.repo.InsertBatch(ctx, *batch)
	if err != nil {
		log.Printf("[Worker %d] Error inserting batch of %d logs: %v", workerID, len(*batch), err)
	} else {
		log.Printf("[Worker %d] Successfully inserted %d logs to MongoDB", workerID, len(*batch))
	}
	*batch = (*batch)[:0]
}
func (p *Pool) Stop() {
	close(p.logChan) 
	p.wg.Wait()      
	log.Println("All workers stopped gracefully")
}
