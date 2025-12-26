package main

import (
	"fmt"
	"gonews/save_service/config"
	"gonews/save_service/internal/bootstrap"
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

	newsStorage := bootstrap.InitPGStorage(cfg)
	_ = bootstrap.InitSaveService(newsStorage, cfg)
}
