package queue

import (
	"context"
	"cushion-isa/internal/models"
)

type QueueClient interface {
	PublishInvestment(inv *models.Investment) error
	StartConsumers(ctx context.Context, handler func(inv *models.Investment) error) error
	Close() error
}
