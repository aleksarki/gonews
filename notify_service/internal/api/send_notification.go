package api

import (
	"context"
	"gonews/protos/pb"
)

func (s *GRPCServer) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	// todo

	return &pb.SendNotificationResponse{Success: true}, nil
}
