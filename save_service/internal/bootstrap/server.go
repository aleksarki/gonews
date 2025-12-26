package bootstrap

import (
	"gonews/protos/pb"
	"gonews/save_service/config"
	"gonews/save_service/internal/api"
	"gonews/save_service/internal/services/saveService"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func InitGRPCServer(saveService *saveService.SaveService, cfg *config.Config) *grpc.Server {
	grpcServer := grpc.NewServer()
	newsServer := api.NewGRPCServer(saveService)
	pb.RegisterSaveServiceServer(grpcServer, newsServer)
	return grpcServer
}

func AppRun(grpcServer *grpc.Server, cfg *config.Config) {
	lis, err := net.Listen("tcp", ":"+string(cfg.GRPC.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Save service gRPC server listening on port %d", cfg.GRPC.Port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		os.Exit(1)
	}
}
