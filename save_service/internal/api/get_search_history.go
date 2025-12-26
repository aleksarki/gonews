package api

import (
	"context"
	"gonews/protos/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) GetSearchHistory(ctx context.Context, req *pb.GetSearchHistoryRequest) (*pb.GetSearchHistoryResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	queries, err := s.saveService.GetSearchHistory(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetSearchHistoryResponse{Queries: queries}, nil
}
