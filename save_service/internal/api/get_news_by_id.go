package api

import (
	"context"
	"gonews/protos/pb"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) GetNewsByIDs(ctx context.Context, req *pb.GetNewsByIDsRequest) (*pb.GetNewsByIDsResponse, error) {
	if len(req.Ids) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no ids provided")
	}

	news, err := s.saveService.GetNewsByIDs(ctx, req.Ids)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbNews := make([]*pb.News, len(news))
	for i, n := range news {
		pbNews[i] = &pb.News{
			Id:          n.ID,
			Source:      n.Source,
			Author:      n.Author,
			Title:       n.Title,
			Description: n.Description,
			Url:         n.URL,
			ImageUrl:    n.ImageURL,
			PublishedAt: n.PublishedAt.Format(time.RFC3339),
		}
	}

	return &pb.GetNewsByIDsResponse{News: pbNews}, nil
}
