package saveService

import (
	"context"
	"gonews/save_service/internal/models"
)

type NewsStorage interface {
	GetNewsByIDs(ctx context.Context, IDs []uint64) ([]*models.News, error)
	UpsertNews(ctx context.Context, news []*models.News) error
}

type SaveService struct {
	newsStorage NewsStorage
}

func NewSaveService(cxt context.Context, newsStorage NewsStorage) *SaveService {
	return &SaveService{newsStorage: newsStorage}
}
