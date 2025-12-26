package api

import (
	"context"
	"gonews/protos/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) AddToSearchHistory(ctx context.Context, req *pb.AddToSearchHistoryRequest) (*pb.AddToSearchHistoryResponse, error) {
	if req.UserId == 0 || req.Query == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and query are required")
	}

	err := s.saveService.AddToSearchHistory(ctx, req.UserId, req.Query, req.Results)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddToSearchHistoryResponse{Success: true}, nil
}
