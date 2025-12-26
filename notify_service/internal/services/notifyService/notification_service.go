package notifyService

import (
	"context"
	"fmt"
	"gonews/notify_service/internal/models"
	"gonews/protos/pb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SaveServiceClient interface {
	GetSubscriptions(ctx context.Context) ([]*models.Subscription, error)
}

type SearchServiceClient interface {
	CheckNewArticles(ctx context.Context, keyword string, lastCheckTime time.Time) ([]*models.News, error)
}

type Producer interface {
	SendNotification(ctx context.Context, userID uint64, keyword string, article models.News) error
	Close() error
}

type NotifyService struct {
	saveClient   pb.SaveServiceClient
	searchClient pb.SearchServiceClient
	producer     Producer
	lastCheck    map[string]time.Time // keyword -> last check time
}

func NewNotifyService(saveServiceAddr, searchServiceAddr string, producer Producer) (*NotifyService, error) {
	// Подключаемся к save service
	saveConn, err := grpc.Dial(saveServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to save service: %w", err)
	}

	// Подключаемся к search service
	searchConn, err := grpc.Dial(searchServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		saveConn.Close()
		return nil, fmt.Errorf("failed to connect to search service: %w", err)
	}

	return &NotifyService{
		saveClient:   pb.NewSaveServiceClient(saveConn),
		searchClient: pb.NewSearchServiceClient(searchConn),
		producer:     producer,
		lastCheck:    make(map[string]time.Time),
	}, nil
}

func (ns *NotifyService) CheckNewArticlesForAllSubscriptions(ctx context.Context) error {
	log.Println("Checking for new articles for all subscriptions...")

	// Получаем все подписки
	resp, err := ns.saveClient.GetSubscriptions(ctx, &pb.GetSubscriptionsRequest{})
	if err != nil {
		return fmt.Errorf("failed to get subscriptions: %w", err)
	}

	for _, sub := range resp.Subscriptions {
		lastCheckTime := ns.lastCheck[sub.Keyword]
		if lastCheckTime.IsZero() {
			lastCheckTime = time.Now().Add(-24 * time.Hour) // Проверяем за последние 24 часа
		}

		err := ns.checkNewArticlesForKeyword(ctx, sub.UserId, sub.Keyword, lastCheckTime)
		if err != nil {
			log.Printf("Error checking keyword '%s' for user %d: %v", sub.Keyword, sub.UserId, err)
			continue
		}

		// Обновляем время последней проверки
		ns.lastCheck[sub.Keyword] = time.Now()
	}

	log.Println("Subscription check completed")
	return nil
}

func (ns *NotifyService) checkNewArticlesForKeyword(ctx context.Context, userID uint64, keyword string, lastCheckTime time.Time) error {
	// Проверяем новые статьи через search service
	resp, err := ns.searchClient.CheckNewArticles(ctx, &pb.CheckNewArticlesRequest{
		Keyword:       keyword,
		LastCheckTime: lastCheckTime.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to check new articles: %w", err)
	}

	// Отправляем уведомления для каждой новой статьи
	for _, article := range resp.NewArticles {
		publishedAt, _ := time.Parse(time.RFC3339, article.PublishedAt)

		news := models.News{
			Source:      article.Source,
			Author:      article.Author,
			Title:       article.Title,
			Description: article.Description,
			URL:         article.Url,
			URLToImage:  article.ImageUrl,
			PublishedAt: publishedAt,
		}

		err := ns.producer.SendNotification(ctx, userID, keyword, news)
		if err != nil {
			log.Printf("Failed to send notification for article '%s': %v", article.Title, err)
			continue
		}

		log.Printf("Notification sent for keyword '%s' to user %d: %s", keyword, userID, article.Title)
	}

	return nil
}

func (ns *NotifyService) Close() error {
	return ns.producer.Close()
}
