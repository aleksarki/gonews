package main

import (
	"context"
	"fmt"
	"gonews/notify_service/config"
	"gonews/notify_service/internal/consumer"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	broker := fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)
	worker := consumer.NewNotificationWorker(
		[]string{broker},
		cfg.Kafka.NotificationTopic,
		cfg.Kafka.ConsumerGroup,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		log.Println("Shutting down worker...")
		cancel()
	}()

	worker.Start(ctx)
	log.Println("Worker stopped")
}
