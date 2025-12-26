package api

import (
	"context"
	"gonews/protos/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) GetFavourites(ctx context.Context, req *pb.GetFavouritesRequest) (*pb.GetFavouritesResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	news, err := s.saveService.GetFavourites(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

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

	return &pb.GetFavouritesResponse{News: protoNews}, nil
}
