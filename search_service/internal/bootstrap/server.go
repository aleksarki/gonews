package bootstrap

import (
	"fmt"
	"gonews/protos/pb"
	"gonews/search_service/config"
	"gonews/search_service/internal/api"
	"gonews/search_service/internal/services/searchService"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func InitGRPCServer(searchService *searchService.SearchService, cfg *config.Config) *grpc.Server {
	grpcServer := grpc.NewServer()
	searchServer := api.NewGRPCServer(searchService)
	pb.RegisterSearchServiceServer(grpcServer, searchServer)
	return grpcServer
}

func AppRun(grpcServer *grpc.Server, cfg *config.Config) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Search service gRPC server listening on port %d", cfg.GRPC.Port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		os.Exit(1)
	}
}
