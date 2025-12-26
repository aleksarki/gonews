package api

import (
	"context"
	"gonews/protos/pb"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) CheckNewArticles(ctx context.Context, req *pb.CheckNewArticlesRequest) (*pb.CheckNewArticlesResponse, error) {
	if req.Keyword == "" {
		return nil, status.Error(codes.InvalidArgument, "keyword is required")
	}

	lastCheckTime := req.LastCheckTime
	if lastCheckTime == "" {
		// Если время не указано, используем время 24 часа назад
		lastCheckTime = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	}

	// Вызываем метод searchService со string
	news, err := s.searchService.CheckNewArticles(ctx, req.Keyword, lastCheckTime)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoNews := make([]*pb.News, len(news))
	for i, n := range news {
		var publishedAtStr string
		if !n.PublishedAt.IsZero() {
			publishedAtStr = n.PublishedAt.Format(time.RFC3339)
		}

		protoNews[i] = &pb.News{
			Id:          n.ID,
			Source:      n.Source,
			Author:      n.Author,
			Title:       n.Title,
			Description: n.Description,
			Url:         n.URL,
			ImageUrl:    n.ImageURL,
			PublishedAt: publishedAtStr,
		}
	}

	return &pb.CheckNewArticlesResponse{
		NewArticles: protoNews,
	}, nil
}
