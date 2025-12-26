package api

import (
	"context"
	"gonews/protos/pb"
)

func (s *GRPCServer) CheckNewArticles(ctx context.Context, req *pb.CheckNewArticlesRequest) (*pb.CheckNewArticlesResponse, error) {
	// todo

	return &pb.CheckNewArticlesResponse{}, nil
}
