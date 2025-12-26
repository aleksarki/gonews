package bootstrap

import (
	"context"
	"gonews/save_service/config"
	"gonews/save_service/internal/services/saveService"
	"gonews/save_service/internal/storage/pgstorage"
)

func InitSaveService(storage *pgstorage.PGStorage, cfg *config.Config) *saveService.SaveService {
	return saveService.NewSaveService(context.Background(), storage)
}
