package queue

import (
	"context"
	"sync"

	"cushion-isa/internal/config"
	"cushion-isa/internal/logger"
	"cushion-isa/internal/models"

	"github.com/sirupsen/logrus"
)

type MemoryQueue struct {
	queue        chan *models.Investment
	numWorkers   int
	workerCancel context.CancelFunc
}

func NewMemoryQueue(cfg *config.QueueConfig) (*MemoryQueue, error) {
	return &MemoryQueue{
		queue:      make(chan *models.Investment, cfg.BufferSize),
		numWorkers: cfg.NumWorkers,
	}, nil
}

func (m *MemoryQueue) PublishInvestment(inv *models.Investment) error {
	m.queue <- inv
	logger.WithFields(logrus.Fields{
		"investment_id": inv.ID,
	}).Info("Investment published to memory queue")
	return nil
}

func (m *MemoryQueue) StartConsumers(ctx context.Context, handler func(inv *models.Investment) error) error {
	ctx, cancel := context.WithCancel(ctx)
	m.workerCancel = cancel

	var wg sync.WaitGroup
	wg.Add(m.numWorkers)

	for i := 0; i < m.numWorkers; i++ {
		workerID := i
		go func() {
			defer wg.Done()
			logger.WithFields(logrus.Fields{
				"worker_id": workerID,
			}).Info("Starting memory queue worker")

			for {
				select {
				case <-ctx.Done():
					return
				case inv := <-m.queue:
					logger.WithFields(logrus.Fields{
						"worker_id":     workerID,
						"investment_id": inv.ID,
					}).Info("Processing investment from memory queue")

					if err := handler(inv); err != nil {
						logger.WithFields(logrus.Fields{
							"worker_id":     workerID,
							"investment_id": inv.ID,
							"error":         err,
						}).Error("Failed to process investment")
					}
				}
			}
		}()
	}

	return nil
}

func (m *MemoryQueue) Close() error {
	if m.workerCancel != nil {
		m.workerCancel()
	}
	close(m.queue)
	return nil
}
