package bootstrap

import (
	"context"
	"fmt"
	"gonews/api_gateway/api"
	"gonews/api_gateway/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func InitHTTPServer(cfg *config.Config) *http.Server {
	handler := api.NewHandler(cfg)

	router := handler.SetupRouter()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server
}

func AppRun(server *http.Server, cfg *config.Config) {
	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запускаем сервер
	go func() {
		log.Printf("API Gateway listening on port %d", cfg.HTTP.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down API Gateway...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("API Gateway stopped")
}
