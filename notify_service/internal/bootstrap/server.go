package bootstrap

import (
	"context"
	"fmt"
	"gonews/notify_service/config"
	"gonews/notify_service/internal/api"
	"gonews/notify_service/internal/services/notifyService"
	"gonews/protos/pb"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func InitGRPCServer(notifyService *notifyService.NotifyService, cfg *config.Config) *grpc.Server {
	grpcServer := grpc.NewServer()
	notificationServer := api.NewGRPCServer(notifyService)
	pb.RegisterNotificationServiceServer(grpcServer, notificationServer)
	return grpcServer
}

func AppRun(grpcServer *grpc.Server, scheduler *notifyService.Scheduler, cfg *config.Config) {
	// Запускаем gRPC сервер
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		log.Printf("Notification service gRPC server listening on port %d", cfg.GRPC.Port)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Запускаем scheduler
	ctx := context.Background()
	scheduler.Start(ctx)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down notification service...")

	// Останавливаем scheduler
	scheduler.Stop()

	// Graceful shutdown gRPC сервера
	grpcServer.GracefulStop()

	log.Println("Notification service stopped")
}
