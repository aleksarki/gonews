package api

import (
	"context"
	"fmt"
	"gonews/api_gateway/config"
	"gonews/protos/pb"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	saveClient         pb.SaveServiceClient
	searchClient       pb.SearchServiceClient
	notificationClient pb.NotificationServiceClient
}

func NewHandler(cfg *config.Config) *Handler {
	// Подключаемся к save service
	saveConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.SaveService.Host, cfg.SaveService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to save service: %v", err))
	}

	// Подключаемся к search service
	searchConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.SearchService.Host, cfg.SearchService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to search service: %v", err))
	}

	// Подключаемся к notification service
	notificationConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.NotifyService.Host, cfg.NotifyService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to notification service: %v", err))
	}

	return &Handler{
		saveClient:         pb.NewSaveServiceClient(saveConn),
		searchClient:       pb.NewSearchServiceClient(searchConn),
		notificationClient: pb.NewNotificationServiceClient(notificationConn),
	}
}

func (h *Handler) SetupRouter() *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := router.Group("/api")
	{
		// User endpoints
		api.POST("/user/create", h.createUser)

		// Search endpoints
		api.GET("/search/news", h.searchNews)
		api.GET("/search/headlines", h.getTopHeadlines)
		api.GET("/search/history/:user_id", h.getSearchHistory)

		// Favourite endpoints
		api.POST("/favourite/set", h.addFavourite)
		api.GET("/favourite/list/:user_id", h.getFavourites)

		// Notification endpoints
		api.POST("/notification/subscribe", h.subscribe)
	}

	// Health check
	router.GET("/health", h.healthCheck)

	return router
}

// Handlers implementation
func (h *Handler) createUser(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.saveClient.CreateUser(c.Request.Context(), &pb.CreateUserRequest{
		Name: req.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": resp.UserId})
}

func (h *Handler) searchNews(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 64)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "valid user_id is required"})
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	// Собираем параметры запроса
	req := &pb.SearchNewsRequest{
		UserId: userID,
		Query:  query,
	}

	// Опциональные строковые параметры
	if sources := c.Query("sources"); sources != "" {
		req.Sources = &sources
	}
	if domains := c.Query("domains"); domains != "" {
		req.Domains = &domains
	}
	if from := c.Query("from"); from != "" {
		req.From = &from
	}
	if to := c.Query("to"); to != "" {
		req.To = &to
	}
	if language := c.Query("language"); language != "" {
		req.Language = &language
	}
	if sortBy := c.Query("sort_by"); sortBy != "" {
		req.SortBy = &sortBy
	}

	// Параметры пагинации (optional int32)
	if pageSize := c.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 {
			pageSizeVal := int32(ps)
			req.PageSize = &pageSizeVal
		}
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageVal := int32(p)
			req.Page = &pageVal
		}
	}

	resp, err := h.searchClient.SearchNews(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"news":          resp.News,
		"total_results": resp.TotalResults,
	})
}

func (h *Handler) getTopHeadlines(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 64)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "valid user_id is required"})
		return
	}

	req := &pb.GetTopHeadlinesRequest{
		UserId: userID,
	}

	// Опциональные строковые параметры
	if country := c.Query("country"); country != "" {
		req.Country = &country
	}
	if category := c.Query("category"); category != "" {
		req.Category = &category
	}
	if sources := c.Query("sources"); sources != "" {
		req.Sources = &sources
	}
	if query := c.Query("q"); query != "" {
		req.Query = &query
	}

	// Параметры пагинации
	if pageSize := c.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 {
			pageSizeVal := int32(ps)
			req.PageSize = &pageSizeVal
		}
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageVal := int32(p)
			req.Page = &pageVal
		}
	}

	resp, err := h.searchClient.GetTopHeadlines(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"news":          resp.News,
		"total_results": resp.TotalResults,
	})
}

func (h *Handler) getSearchHistory(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "valid user_id is required"})
		return
	}

	resp, err := h.saveClient.GetSearchHistory(c.Request.Context(), &pb.GetSearchHistoryRequest{
		UserId: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": resp.Queries})
}

func (h *Handler) addFavourite(c *gin.Context) {
	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
		NewsID uint64 `json:"news_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.saveClient.AddFavourite(c.Request.Context(), &pb.AddFavouriteRequest{
		UserId: req.UserID,
		NewsId: req.NewsID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": resp.Success})
}

func (h *Handler) getFavourites(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "valid user_id is required"})
		return
	}

	resp, err := h.saveClient.GetFavourites(c.Request.Context(), &pb.GetFavouritesRequest{
		UserId: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"favourites": resp.News})
}

func (h *Handler) subscribe(c *gin.Context) {
	var req struct {
		UserID  uint64 `json:"user_id" binding:"required"`
		Keyword string `json:"keyword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.saveClient.Subscribe(c.Request.Context(), &pb.SubscribeRequest{
		UserId:  req.UserID,
		Keyword: req.Keyword,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": resp.Success})
}

func (h *Handler) healthCheck(c *gin.Context) {
	// Проверяем соединения со всеми сервисами
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	health := map[string]string{
		"api_gateway": "ok",
	}

	// Проверка save service
	_, err := h.saveClient.CreateUser(ctx, &pb.CreateUserRequest{Name: "health_check"})
	if err != nil {
		// Если ошибка "user already exists" - это нормально
		if err.Error() != "user already exists" {
			health["save_service"] = "error: " + err.Error()
		} else {
			health["save_service"] = "ok"
		}
	} else {
		health["save_service"] = "ok"
	}

	// Проверка search service
	_, err = h.searchClient.SearchNews(ctx, &pb.SearchNewsRequest{
		UserId: 1,
		Query:  "test",
	})
	if err != nil {
		health["search_service"] = "error: " + err.Error()
	} else {
		health["search_service"] = "ok"
	}

	// Проверка notification service
	_, err = h.notificationClient.SendNotification(ctx, &pb.SendNotificationRequest{
		UserId:  1,
		Message: "health check",
	})
	if err != nil {
		health["notification_service"] = "error: " + err.Error()
	} else {
		health["notification_service"] = "ok"
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "services": health})
}
