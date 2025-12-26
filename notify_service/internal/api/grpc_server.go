package api

import (
	"gonews/protos/pb"
)

type NotifyService interface {
	/// fixme
}

type GRPCServer struct {
	pb.UnimplementedNotificationServiceServer
	notifyService NotifyService
}

func NewGRPCServer(notifyService NotifyService) *GRPCServer {
	return &GRPCServer{
		notifyService: notifyService,
	}
}
