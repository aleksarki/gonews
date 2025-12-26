package bootstrap

import (
	"fmt"
	"gonews/notify_service/config"
	"gonews/notify_service/internal/producer"
	"gonews/notify_service/internal/services/notifyService"
	"time"
)

func InitNotifyService(producer *producer.KafkaProducer, cfg *config.Config) (*notifyService.NotifyService, error) {
	saveServiceAddr := fmt.Sprintf("%s:%d", cfg.SaveService.Host, cfg.SaveService.Port)
	searchServiceAddr := fmt.Sprintf("%s:%d", cfg.SearchService.Host, cfg.SearchService.Port)

	return notifyService.NewNotifyService(saveServiceAddr, searchServiceAddr, producer)
}

func InitScheduler(notifyService_ *notifyService.NotifyService, cfg *config.Config) *notifyService.Scheduler {
	interval := time.Duration(cfg.Scheduler.CheckIntervalMinutes) * time.Minute
	return notifyService.NewScheduler(notifyService_, interval)
}
