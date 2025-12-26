package main

import (
	"fmt"
	"gonews/notify_service/config"
	"gonews/notify_service/internal/bootstrap"
	"os"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("Config load error: %v", err))
	}

	kafkaProducer := bootstrap.InitKafkaProducer(cfg)
	notifyService, err := bootstrap.InitNotifyService(kafkaProducer, cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize notify service: %v", err))
	}

	scheduler := bootstrap.InitScheduler(notifyService, cfg)
	grpcServer := bootstrap.InitGRPCServer(notifyService, cfg)

	bootstrap.AppRun(grpcServer, scheduler, cfg)
}
