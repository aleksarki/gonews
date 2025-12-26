package api

import (
	"context"
	"gonews/protos/pb"
	"gonews/search_service/internal/services/searchService"

	"google.golang.org/grpc"
)

type SearchService interface {
	CheckNewArticles(ctx context.Context, keyword string, lastCheckTime string) ([]*searchService.News, error)
	GetTopHeadlines(ctx context.Context, req *searchService.TopHeadlinesRequest) ([]*searchService.News, int, error)
	SearchNews(ctx context.Context, req *searchService.SearchRequest) ([]*searchService.News, int, error)
}

type GRPCServer struct {
	pb.UnimplementedSearchServiceServer
	searchService SearchService
}

func NewGRPCServer(searchService SearchService) *GRPCServer {
	return &GRPCServer{
		searchService: searchService,
	}
}

func (s *GRPCServer) Register(server *grpc.Server) {
	pb.RegisterSearchServiceServer(server, s)
}
