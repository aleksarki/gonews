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
	saveConn     *grpc.ClientConn
	searchConn   *grpc.ClientConn
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
		saveConn:     saveConn,
		searchConn:   searchConn,
		lastCheck:    make(map[string]time.Time),
	}, nil
}

// GetSaveClient возвращает клиент для save service
func (ns *NotifyService) GetSaveClient() pb.SaveServiceClient {
	return ns.saveClient
}

// GetSearchClient возвращает клиент для search service
func (ns *NotifyService) GetSearchClient() pb.SearchServiceClient {
	return ns.searchClient
}

// CheckNewArticlesForAllSubscriptions - проверяет новые статьи для всех подписок
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

		err := ns.CheckNewArticlesForKeyword(ctx, sub.UserId, sub.Keyword, lastCheckTime)
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

func (ns *NotifyService) CheckNewArticlesForKeyword(ctx context.Context, userID uint64, keyword string, lastCheckTime time.Time) error {
	// ignore result
	_, err := ns.CheckNewArticlesForSubscription(ctx, userID, keyword, lastCheckTime)
	return err
}

// закрыть соединения
func (ns *NotifyService) Close() error {
	var errs []error

	if ns.saveConn != nil {
		if err := ns.saveConn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if ns.searchConn != nil {
		if err := ns.searchConn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if ns.producer != nil {
		if err := ns.producer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}

	return nil
}

// GetSubscriptionsByKeyword - получает все подписки по ключевому слову
func (ns *NotifyService) GetSubscriptionsByKeyword(ctx context.Context, keyword string) ([]*models.Subscription, error) {
	resp, err := ns.saveClient.GetSubscriptions(ctx, &pb.GetSubscriptionsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	var subscriptions []*models.Subscription
	for _, sub := range resp.Subscriptions {
		if sub.Keyword == keyword {
			subscriptions = append(subscriptions, &models.Subscription{
				ID:      sub.Id,
				UserID:  sub.UserId,
				Keyword: sub.Keyword,
			})
		}
	}

	return subscriptions, nil
}

// CheckNewArticlesForSubscription - публичный метод для проверки новых статей
func (ns *NotifyService) CheckNewArticlesForSubscription(ctx context.Context, userID uint64, keyword string, lastCheckTime time.Time) ([]*models.News, error) {
	log.Printf("Checking new articles for user %d, keyword: %s, since: %v", userID, keyword, lastCheckTime)

	// Проверяем новые статьи через search service
	resp, err := ns.searchClient.CheckNewArticles(ctx, &pb.CheckNewArticlesRequest{
		Keyword:       keyword,
		LastCheckTime: lastCheckTime.Format(time.RFC3339),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check new articles: %w", err)
	}

	// Конвертируем в модели
	var news []*models.News
	for _, article := range resp.NewArticles {
		var publishedAt time.Time
		if article.PublishedAt != "" {
			publishedAt, _ = time.Parse(time.RFC3339, article.PublishedAt)
		}

		newsItem := &models.News{
			Source:      article.Source,
			Author:      article.Author,
			Title:       article.Title,
			Description: article.Description,
			URL:         article.Url,
			URLToImage:  article.ImageUrl,
			PublishedAt: publishedAt,
		}

		news = append(news, newsItem)

		// Отправляем уведомление
		err := ns.producer.SendNotification(ctx, userID, keyword, *newsItem)
		if err != nil {
			log.Printf("Failed to send notification for article '%s': %v", article.Title, err)
			continue
		}

		log.Printf("Notification sent for keyword '%s' to user %d: %s", keyword, userID, article.Title)
	}

	log.Printf("Found %d new articles for user %d, keyword: %s", len(news), userID, keyword)
	return news, nil
}

// SendNotification - отправляет уведомление
func (ns *NotifyService) SendNotification(ctx context.Context, userID uint64, keyword string, article models.News) error {
	return ns.producer.SendNotification(ctx, userID, keyword, article)
}
