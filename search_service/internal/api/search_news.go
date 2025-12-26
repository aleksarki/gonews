package api

import (
	"context"
	"gonews/protos/pb"
	"gonews/search_service/internal/services/searchService"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) SearchNews(ctx context.Context, req *pb.SearchNewsRequest) (*pb.SearchNewsResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	if req.Query == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	// Convert gRPC request to service request
	searchReq := &searchService.SearchRequest{
		UserID: req.UserId,
		Query:  req.Query,
	}

	// опциональные поля
	if req.Sources != nil {
		searchReq.Sources = *req.Sources
	}
	if req.Domains != nil {
		searchReq.Domains = *req.Domains
	}
	if req.From != nil {
		searchReq.From = *req.From
	}
	if req.To != nil {
		searchReq.To = *req.To
	}
	if req.Language != nil {
		searchReq.Language = *req.Language
	}
	if req.SortBy != nil {
		searchReq.SortBy = *req.SortBy
	}
	if req.PageSize != nil {
		searchReq.PageSize = int(*req.PageSize)
	}
	if req.Page != nil {
		searchReq.Page = int(*req.Page)
	}

	news, totalResults, err := s.searchService.SearchNews(ctx, searchReq)
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

	return &pb.SearchNewsResponse{
		News:         protoNews,
		TotalResults: int32(totalResults),
	}, nil
}
