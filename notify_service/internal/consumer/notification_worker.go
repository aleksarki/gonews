package consumer

import (
	"context"
	"encoding/json"
	"gonews/notify_service/internal/models"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type NotificationWorker struct {
	reader *kafka.Reader
}

func NewNotificationWorker(brokers []string, topic, groupID string) *NotificationWorker {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	})

	return &NotificationWorker{
		reader: reader,
	}
}

func (w *NotificationWorker) Start(ctx context.Context) {
	log.Println("Starting notification worker...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping notification worker...")
			w.reader.Close()
			return
		default:
			msg, err := w.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			w.processMessage(msg.Value)
		}
	}
}

func (w *NotificationWorker) processMessage(message []byte) {
	var notification models.NotificationMessage
	err := json.Unmarshal(message, &notification)
	if err != nil {
		log.Printf("Failed to unmarshal notification: %v", err)
		return
	}

	// fake sending notifications
	log.Printf("=== NOTIFICATION SENT ===")
	log.Printf("Event ID: %s", notification.EventID)
	log.Printf("To topic: %s", notification.NotifTopic)
	log.Printf("Article: %s", notification.Article.Title)
	log.Printf("Author: %s", notification.Article.Author)
	log.Printf("Published: %s", notification.Article.PublishedAt.Format("2006-01-02 15:04"))
	log.Printf("URL: %s", notification.Article.URL)
	log.Printf("Timestamp: %s", notification.Timestamp.Format("2006-01-02 15:04:05"))
	log.Printf("=========================\n")
}

func (w *NotificationWorker) Close() error {
	return w.reader.Close()
}
