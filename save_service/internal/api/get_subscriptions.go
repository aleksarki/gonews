package api

import (
	"context"
	"gonews/protos/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) GetSubscriptions(ctx context.Context, req *pb.GetSubscriptionsRequest) (*pb.GetSubscriptionsResponse, error) {
	subscriptions, err := s.saveService.GetSubscriptions(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoSubs := make([]*pb.Subscription, len(subscriptions))
	for i, sub := range subscriptions {
		protoSubs[i] = &pb.Subscription{
			Id:      sub.ID,
			UserId:  sub.UserID,
			Keyword: sub.Keyword,
		}
	}

	return &pb.GetSubscriptionsResponse{Subscriptions: protoSubs}, nil
}
