package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"gonews/notify_service/internal/models"
	"time"

	"github.com/samber/lo"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
	topic  string
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}

	return &KafkaProducer{
		writer: writer,
		topic:  topic,
	}
}

func (kp *KafkaProducer) SendNotification(ctx context.Context, userID uint64, keyword string, article models.News) error {
	message := models.NotificationMessage{
		EventID:    lo.RandomString(10, lo.LettersCharset),
		EventType:  "notification",
		NotifTopic: keyword,
		Article:    article,
		Timestamp:  time.Now(),
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal notification message: %w", err)
	}

	// Отправляем сообщение в Kafka
	err = kp.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(fmt.Sprintf("user_%d", userID)),
			Value: messageJSON,
			Time:  time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}
