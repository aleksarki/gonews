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

func (s *SearchService) CheckNewArticles(ctx context.Context, keyword, lastCheckTime string) ([]*News, error) {
	// Implementation for notification service
	// todo
	return nil, nil
}
