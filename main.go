package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"cushion-isa/internal/api"
	"cushion-isa/internal/config"
	"cushion-isa/internal/db"
	"cushion-isa/internal/logger"
	"cushion-isa/internal/queue"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	logger.InitLogger(cfg.Logger.Level)

	// Initialize database connection
	dbClient, err := db.NewClient(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to create database client:", err)
	}
	defer dbClient.Close()

	// Run database migrations
	if err := dbClient.Migrate(); err != nil {
		logger.Fatal("Failed to run database migrations:", err)
	}

	// Initialize queue client based on configuration
	var queueClient queue.QueueClient
	switch cfg.Queue.Type {
	case "memory":
		queueClient, err = queue.NewMemoryQueue(&cfg.Queue)
	case "kafka":
		queueClient, err = queue.NewKafkaClient(&cfg.Queue)
	default:
		logger.Fatalf("Unsupported queue type: %s", cfg.Queue.Type)
	}
	if err != nil {
		logger.Fatal("Failed to create queue client:", err)
	}
	defer queueClient.Close()

	// Initialize API handler
	handler := api.NewHandler(dbClient, queueClient)

	// Start HTTP server
	server := api.NewServer(cfg.Server, handler)
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal("Failed to start server:", err)
		}
	}()

	// Start queue consumers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := queueClient.StartConsumers(ctx, handler.ProcessInvestment); err != nil {
		logger.Fatal("Failed to start consumers:", err)
	}

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down gracefully...")
	if err := server.Stop(ctx); err != nil {
		logger.Error("Error during server shutdown:", err)
	}
}
