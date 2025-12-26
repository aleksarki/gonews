package main

import (
	"fmt"
	"gonews/api_gateway/bootstrap"
	"gonews/api_gateway/config"
	"os"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("Config load error: %v", err))
	}

	server := bootstrap.InitHTTPServer(cfg)

	bootstrap.AppRun(server, cfg)
}
