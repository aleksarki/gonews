package main

import (
	"fmt"
	"gonews/search_service/config"
	"os"
)

func main() {
	path := os.Getenv("configPath")
	if path == "" {
		path = "/config/config.yaml"
	}
	_, err := config.LoadConfig(path)
	if err != nil {
		panic(fmt.Sprintf("config load error: %v", err))
	}
}
