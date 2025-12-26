package api

import (
	"context"
	"gonews/protos/pb"
	"gonews/save_service/internal/models"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) SaveNews(ctx context.Context, req *pb.SaveNewsRequest) (*pb.SaveNewsResponse, error) {
	if len(req.News) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no news to save")
	}

	news := make([]*models.News, len(req.News))
	for i, n := range req.News {
		publishedAt, _ := time.Parse(time.RFC3339, n.PublishedAt)

		news[i] = &models.News{
			Source:      n.Source,
			Author:      n.Author,
			Title:       n.Title,
			Description: n.Description,
			URL:         n.Url,
			ImageURL:    n.ImageUrl,
			PublishedAt: publishedAt,
		}
	}

	err := s.saveService.SaveNews(ctx, news)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SaveNewsResponse{Success: true}, nil
}
