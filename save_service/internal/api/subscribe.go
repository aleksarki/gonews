package api

import (
	"context"
	"gonews/protos/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) Subscribe(ctx context.Context, req *pb.SubscribeRequest) (*pb.SubscribeResponse, error) {
	if req.UserId == 0 || req.Keyword == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and keyword are required")
	}

	err := s.saveService.Subscribe(ctx, req.UserId, req.Keyword)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SubscribeResponse{Success: true}, nil
}
