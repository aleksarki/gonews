package api

import (
	"context"
	"gonews/notify_service/internal/models"
	"gonews/protos/pb"
	"time"
)

type NotifyService interface {
	CheckNewArticlesForAllSubscriptions(ctx context.Context) error
	CheckNewArticlesForKeyword(ctx context.Context, userID uint64, keyword string, lastCheckTime time.Time) error
	Close() error
	GetSubscriptionsByKeyword(ctx context.Context, keyword string) ([]*models.Subscription, error)
	CheckNewArticlesForSubscription(ctx context.Context, userID uint64, keyword string, lastCheckTime time.Time) ([]*models.News, error)
	SendNotification(ctx context.Context, userID uint64, keyword string, article models.News) error
	GetSaveClient() pb.SaveServiceClient
	GetSearchClient() pb.SearchServiceClient
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
