package api

import (
	"context"
	"gonews/protos/pb"
	"gonews/save_service/internal/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SaveService interface {
	CreateUser(ctx context.Context, name string) (uint64, error)
	AddFavourite(ctx context.Context, userID, newsID uint64) error
	GetFavourites(ctx context.Context, userID uint64) ([]*models.News, error)
	AddToSearchHistory(ctx context.Context, userID uint64, query string, results []uint64) error
	GetSearchHistory(ctx context.Context, userID uint64) ([]string, error)
	Subscribe(ctx context.Context, userID uint64, keyword string) error
	GetSubscriptions(ctx context.Context) ([]*models.Subscription, error)
	MarkNewsAsSeen(ctx context.Context, userID, newsID uint64) error
	SaveNews(ctx context.Context, news []*models.News) error
	GetNewsByIDs(ctx context.Context, IDs []uint64) ([]*models.News, error)
}

type GRPCServer struct {
	pb.UnimplementedSaveServiceServer
	saveService SaveService
}

func NewGRPCServer(saveService SaveService) *GRPCServer {
	return &GRPCServer{
		saveService: saveService,
	}
}

func (s *GRPCServer) AddFavourite(ctx context.Context, req *pb.AddFavouriteRequest) (*pb.AddFavouriteResponse, error) {
	if req.UserId == 0 || req.NewsId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id and news_id are required")
	}

	err := s.saveService.AddFavourite(ctx, req.UserId, req.NewsId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddFavouriteResponse{Success: true}, nil
}

func (s *GRPCServer) Register(server *grpc.Server) {
	pb.RegisterSaveServiceServer(server, s)
}
