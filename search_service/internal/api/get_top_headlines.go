package api

import (
	"context"
	"gonews/protos/pb"
	"gonews/search_service/internal/services/searchService"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) GetTopHeadlines(ctx context.Context, req *pb.GetTopHeadlinesRequest) (*pb.GetTopHeadlinesResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Convert gRPC request to service request
	headlinesReq := &searchService.TopHeadlinesRequest{
		UserID: req.UserId,
	}

	// опциональные поля
	if req.Country != nil {
		headlinesReq.Country = *req.Country
	}
	if req.Category != nil {
		headlinesReq.Category = *req.Category
	}
	if req.Sources != nil {
		headlinesReq.Sources = *req.Sources
	}
	if req.Query != nil {
		headlinesReq.Query = *req.Query
	}
	if req.PageSize != nil {
		headlinesReq.PageSize = int(*req.PageSize)
	}
	if req.Page != nil {
		headlinesReq.Page = int(*req.Page)
	}

	news, totalResults, err := s.searchService.GetTopHeadlines(ctx, headlinesReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to protobuf response
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

	return &pb.GetTopHeadlinesResponse{
		News:         protoNews,
		TotalResults: int32(totalResults),
	}, nil
}
