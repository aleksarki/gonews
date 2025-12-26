package bootstrap

import (
	"gonews/search_service/config"
	"gonews/search_service/internal/newsapi"
	"log"
)

func InitNewsAPIClient(cfg *config.Config) *newsapi.Client {
	if cfg.NewsAPI.APIKey == "" {
		log.Fatal("NewsAPI API key is required")
	}
	return newsapi.NewClient(cfg.NewsAPI.APIKey)
}
