package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cushion-isa/internal/config"
	"cushion-isa/internal/logger"
	"cushion-isa/internal/models"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type KafkaClient struct {
	producer     sarama.SyncProducer
	consumer     sarama.ConsumerGroup
	topic        string
	numWorkers   int
	workerCancel context.CancelFunc
}

func NewKafkaClient(cfg *config.QueueConfig) (*KafkaClient, error) {
	// Producer config
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5
	producerConfig.Producer.Return.Successes = true

	// Create producer
	producer, err := sarama.NewSyncProducer(cfg.Brokers, producerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	// Consumer config
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Create consumer group
	consumer, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.ConsumerGroup, consumerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &KafkaClient{
		producer:   producer,
		consumer:   consumer,
		topic:      cfg.Topic,
		numWorkers: cfg.NumWorkers,
	}, nil
}

func (k *KafkaClient) PublishInvestment(inv *models.Investment) error {
	payload, err := json.Marshal(inv)
	if err != nil {
		return fmt.Errorf("failed to marshal investment: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.StringEncoder(payload),
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"partition":     partition,
		"offset":        offset,
		"investment_id": inv.ID,
	}).Info("Investment message published to Kafka")

	return nil
}

func (k *KafkaClient) StartConsumers(ctx context.Context, handler func(inv *models.Investment) error) error {
	ctx, cancel := context.WithCancel(ctx)
	k.workerCancel = cancel

	var wg sync.WaitGroup
	wg.Add(k.numWorkers)

	// Start consumer workers
	for i := 0; i < k.numWorkers; i++ {
		workerID := i
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					if err := k.consume(ctx, handler, workerID); err != nil {
						logger.WithFields(logrus.Fields{
							"worker_id": workerID,
							"error":     err,
						}).Error("Consumer error")
					}
				}
			}
		}()
	}

	go func() {
		<-ctx.Done()
		if err := k.consumer.Close(); err != nil {
			logger.Errorf("Error closing consumer: %v", err)
		}
	}()

	return nil
}

func (k *KafkaClient) consume(ctx context.Context, handler func(inv *models.Investment) error, workerID int) error {
	// Consume messages
	topics := []string{k.topic}
	groupHandler := &ConsumerGroupHandler{
		handler:  handler,
		workerID: workerID,
	}

	return k.consumer.Consume(ctx, topics, groupHandler)
}

type ConsumerGroupHandler struct {
	handler  func(inv *models.Investment) error
	workerID int
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			var investment models.Investment
			if err := json.Unmarshal(message.Value, &investment); err != nil {
				logger.WithFields(logrus.Fields{
					"worker_id": h.workerID,
					"error":     err,
				}).Error("Failed to unmarshal investment message")
				continue
			}

			logger.WithFields(logrus.Fields{
				"worker_id":     h.workerID,
				"investment_id": investment.ID,
			}).Info("Processing investment message")

			if err := h.handler(&investment); err != nil {
				logger.WithFields(logrus.Fields{
					"worker_id":     h.workerID,
					"investment_id": investment.ID,
					"error":         err,
				}).Error("Failed to process investment")
				continue
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (k *KafkaClient) Close() error {
	if k.workerCancel != nil {
		k.workerCancel()
	}

	if err := k.producer.Close(); err != nil {
		return fmt.Errorf("failed to close producer: %w", err)
	}

	if err := k.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %w", err)
	}

	return nil
}
