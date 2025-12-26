package main

import (
	"fmt"
	"gonews/search_service/config"
	"gonews/search_service/internal/bootstrap"
	"os"
)

func main() {
	path := os.Getenv("configPath")
	if path == "" {
		path = "/config/config.yaml"
	}
	cfg, err := config.LoadConfig(path)
	if err != nil {
		panic(fmt.Sprintf("config load error: %v", err))
	}

	redisStorage := bootstrap.InitRedisStorage(cfg)
	newsAPIClient := bootstrap.InitNewsAPIClient(cfg)
	searchService := bootstrap.InitSearchService(newsAPIClient, redisStorage, cfg)
	grpcServer := bootstrap.InitGRPCServer(searchService, cfg)
	bootstrap.AppRun(grpcServer, cfg)
}
