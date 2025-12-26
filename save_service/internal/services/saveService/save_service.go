package saveService

import (
	"context"
	"gonews/save_service/internal/models"
)

type NewsStorage interface {
	GetNewsByIDs(ctx context.Context, IDs []uint64) ([]*models.News, error)
	UpsertNews(ctx context.Context, news []*models.News) error
	CreateUser(ctx context.Context, name string) (uint64, error)
	AddFavourite(ctx context.Context, userID, newsID uint64) error
	GetFavourites(ctx context.Context, userID uint64) ([]*models.News, error)
	AddToSearchHistory(ctx context.Context, userID uint64, query string, results []uint64) error
	GetSearchHistory(ctx context.Context, userID uint64) ([]string, error)
	Subscribe(ctx context.Context, userID uint64, keyword string) error
	GetSubscriptions(ctx context.Context) ([]*models.Subscription, error)
	MarkNewsAsSeen(ctx context.Context, userID, newsID uint64) error
}

type SaveService struct {
	newsStorage NewsStorage
}

func NewSaveService(ctx context.Context, newsStorage NewsStorage) *SaveService {
	return &SaveService{newsStorage: newsStorage}
}

func (s *SaveService) CreateUser(ctx context.Context, name string) (uint64, error) {
	return s.newsStorage.CreateUser(ctx, name)
}

func (s *SaveService) AddFavourite(ctx context.Context, userID, newsID uint64) error {
	return s.newsStorage.AddFavourite(ctx, userID, newsID)
}

func (s *SaveService) GetFavourites(ctx context.Context, userID uint64) ([]*models.News, error) {
	return s.newsStorage.GetFavourites(ctx, userID)
}

func (s *SaveService) AddToSearchHistory(ctx context.Context, userID uint64, query string, results []uint64) error {
	return s.newsStorage.AddToSearchHistory(ctx, userID, query, results)
}

func (s *SaveService) GetSearchHistory(ctx context.Context, userID uint64) ([]string, error) {
	return s.newsStorage.GetSearchHistory(ctx, userID)
}

func (s *SaveService) Subscribe(ctx context.Context, userID uint64, keyword string) error {
	return s.newsStorage.Subscribe(ctx, userID, keyword)
}

func (s *SaveService) GetSubscriptions(ctx context.Context) ([]*models.Subscription, error) {
	return s.newsStorage.GetSubscriptions(ctx)
}

func (s *SaveService) MarkNewsAsSeen(ctx context.Context, userID, newsID uint64) error {
	return s.newsStorage.MarkNewsAsSeen(ctx, userID, newsID)
}

func (s *SaveService) SaveNews(ctx context.Context, news []*models.News) error {
	return s.newsStorage.UpsertNews(ctx, news)
}

func (s *SaveService) GetNewsByIDs(ctx context.Context, IDs []uint64) ([]*models.News, error) {
	return s.newsStorage.GetNewsByIDs(ctx, IDs)
}
