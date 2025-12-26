package api

import (
	"context"
	"gonews/protos/pb"
	"gonews/search_service/internal/services/searchService"

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
		UserID:   req.UserId,
		Query:    req.Query,
		Sources:  *req.Sources,
		Domains:  *req.Domains,
		From:     *req.From,
		To:       *req.To,
		Language: *req.Language,
		SortBy:   *req.SortBy,
		PageSize: int(*req.PageSize),
		Page:     int(*req.Page),
	}

	news, totalResults, err := s.searchService.SearchNews(ctx, searchReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to protobuf response
	protoNews := make([]*pb.News, len(news))
	for i, n := range news {
		protoNews[i] = &pb.News{
			Id:          n.ID,
			Source:      n.Source,
			Author:      n.Author,
			Title:       n.Title,
			Description: n.Description,
			Url:         n.URL,
			ImageUrl:    n.ImageURL,
			PublishedAt: n.PublishedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return &pb.SearchNewsResponse{
		News:         protoNews,
		TotalResults: int32(totalResults),
	}, nil
}
