package searchService

import (
	"context"
	"encoding/json"
	"fmt"
	"gonews/protos/pb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type News struct {
	ID          uint64
	Source      string
	Author      string
	Title       string
	Description string
	URL         string
	ImageURL    string
	PublishedAt time.Time
	Content     string
}

type SearchRequest struct {
	UserID   uint64
	Query    string
	Sources  string
	Domains  string
	From     string
	To       string
	Language string
	SortBy   string
	PageSize int
	Page     int
}

type TopHeadlinesRequest struct {
	UserID   uint64
	Country  string
	Category string
	Sources  string
	Query    string
	PageSize int
	Page     int
}

type NewsAPIClient interface {
	SearchEverything(ctx context.Context, req *SearchRequest) ([]*News, int, error)
	GetTopHeadlines(ctx context.Context, req *TopHeadlinesRequest) ([]*News, int, error)
	CheckNewArticles(ctx context.Context, keyword, lastCheckTimeStr string) ([]*News, error)
}

type CacheStorage interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}

type SearchService struct {
	newsAPI         NewsAPIClient
	cache           CacheStorage
	saveServiceAddr string
	newsAPIKey      string
}

func NewSearchService(newsAPI NewsAPIClient, cache CacheStorage, saveServiceAddr, newsAPIKey string) *SearchService {
	return &SearchService{
		newsAPI:         newsAPI,
		cache:           cache,
		saveServiceAddr: saveServiceAddr,
		newsAPIKey:      newsAPIKey,
	}
}

func (s *SearchService) SearchNews(ctx context.Context, req *SearchRequest) ([]*News, int, error) {
	// Check cache
	cacheKey := fmt.Sprintf("search:%s:%s:%s:%d:%d",
		req.Query, req.From, req.To, req.PageSize, req.Page)

	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var cachedResult struct {
			News  []*News `json:"news"`
			Total int     `json:"total"`
		}
		if err := json.Unmarshal([]byte(cached), &cachedResult); err == nil {
			return cachedResult.News, cachedResult.Total, nil
		}
	}

	// Call external API
	news, total, err := s.newsAPI.SearchEverything(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	// Save to cache
	resultJSON, _ := json.Marshal(map[string]interface{}{
		"news":  news,
		"total": total,
	})
	s.cache.Set(ctx, cacheKey, string(resultJSON), 10*time.Minute)

	// Save to database via gRPC
	conn, err := grpc.Dial(s.saveServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err == nil {
		defer conn.Close()

		client := pb.NewSaveServiceClient(conn)

		// Convert to protobuf
		pbNews := make([]*pb.News, len(news))
		for i, n := range news {
			pbNews[i] = &pb.News{
				Source:      n.Source,
				Author:      n.Author,
				Title:       n.Title,
				Description: n.Description,
				Url:         n.URL,
				ImageUrl:    n.ImageURL,
				PublishedAt: n.PublishedAt.Format(time.RFC3339),
			}
		}

		// Save news
		client.SaveNews(ctx, &pb.SaveNewsRequest{News: pbNews})

		// Save search history
		newsIDs := make([]uint64, len(news))
		for i, n := range news {
			newsIDs[i] = n.ID
		}
		client.AddToSearchHistory(ctx, &pb.AddToSearchHistoryRequest{
			UserId:  req.UserID,
			Query:   req.Query,
			Results: newsIDs,
		})
	}

	return news, total, nil
}

func (s *SearchService) GetTopHeadlines(ctx context.Context, req *TopHeadlinesRequest) ([]*News, int, error) {
	// Similar implementation...
	// todo
	return nil, 0, nil
}

func (s *SearchService) CheckNewArticles(ctx context.Context, keyword, lastCheckTimeStr string) ([]*News, error) {
	// Парсим строку времени в time.Time
	var lastCheckTime time.Time
	var err error

	if lastCheckTimeStr != "" {
		lastCheckTime, err = time.Parse(time.RFC3339, lastCheckTimeStr)
		if err != nil {
			// Если не удалось распарсить, используем время по умолчанию (24 часа назад)
			lastCheckTime = time.Now().Add(-24 * time.Hour)
		}
	} else {
		// Если строка пустая, используем время по умолчанию
		lastCheckTime = time.Now().Add(-24 * time.Hour)
	}

	// Создаем запрос для поиска статей с момента lastCheckTime
	fromTime := lastCheckTime.Format("2006-01-02T15:04:05Z")

	searchReq := &SearchRequest{
		Query:    keyword,
		From:     fromTime,
		SortBy:   "publishedAt", // Сортируем по дате публикации
		PageSize: 50,            // Проверяем до 50 новых статей
	}

	// Ищем новые статьи
	news, totalResults, err := s.newsAPI.SearchEverything(ctx, searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to check new articles: %w", err)
	}

	if totalResults == 0 {
		return []*News{}, nil // Нет новых статей
	}

	// Сохраняем в кэш
	cacheKey := fmt.Sprintf("check:%s:%s", keyword, fromTime)
	newsJSON, _ := json.Marshal(news)
	s.cache.Set(ctx, cacheKey, string(newsJSON), 30*time.Minute)

	// Сохраняем в базу через save service
	conn, err := grpc.Dial(s.saveServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err == nil {
		defer conn.Close()

		client := pb.NewSaveServiceClient(conn)

		// Конвертируем в protobuf
		pbNews := make([]*pb.News, len(news))
		for i, n := range news {
			var publishedAtStr string
			if !n.PublishedAt.IsZero() {
				publishedAtStr = n.PublishedAt.Format(time.RFC3339)
			}

			pbNews[i] = &pb.News{
				Source:      n.Source,
				Author:      n.Author,
				Title:       n.Title,
				Description: n.Description,
				Url:         n.URL,
				ImageUrl:    n.ImageURL,
				PublishedAt: publishedAtStr,
			}
		}

		// Сохраняем новости
		client.SaveNews(ctx, &pb.SaveNewsRequest{News: pbNews})
	}

	return news, nil
}
