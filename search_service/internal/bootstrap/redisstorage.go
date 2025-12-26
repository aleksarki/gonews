package bootstrap

import (
	"fmt"
	"gonews/search_service/config"
	"gonews/search_service/internal/storage"
	"log"
)

func InitRedisStorage(cfg *config.Config) *storage.RedisStorage {
	redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	redisStorage, err := storage.NewRedisStorage(redisAddr)
	if err != nil {
		log.Fatalf("Failed to initialize Redis storage: %v", err)
	}
	return redisStorage
}
