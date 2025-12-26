package bootstrap

import (
	"fmt"
	"gonews/search_service/config"
	"gonews/search_service/internal/newsapi"
	"gonews/search_service/internal/services/searchService"
	"gonews/search_service/internal/storage"
)

func InitSearchService(newsAPI *newsapi.Client, cache *storage.RedisStorage, cfg *config.Config) *searchService.SearchService {
	saveServiceAddr := fmt.Sprintf("%s:%d", cfg.SaveService.Host, cfg.SaveService.Port)
	return searchService.NewSearchService(newsAPI, cache, saveServiceAddr, cfg.NewsAPI.APIKey)
}
