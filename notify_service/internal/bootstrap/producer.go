package bootstrap

import (
	"fmt"
	"gonews/notify_service/config"
	"gonews/notify_service/internal/producer"
)

func InitKafkaProducer(cfg *config.Config) *producer.KafkaProducer {
	broker := fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)
	return producer.NewKafkaProducer([]string{broker}, cfg.Kafka.NotificationTopic)
}
