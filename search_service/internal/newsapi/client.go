package newsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gonews/search_service/internal/services/searchService"
)

const (
	everythingURL   = "https://newsapi.org/v2/everything"
	topHeadlinesURL = "https://newsapi.org/v2/top-headlines"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Response from NewsAPI
type NewsAPIResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Article struct {
	Source      Source `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	URLToImage  string `json:"urlToImage"`
	PublishedAt string `json:"publishedAt"`
	Content     string `json:"content"`
}

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) SearchEverything(ctx context.Context, req *searchService.SearchRequest) ([]*searchService.News, int, error) {
	// Build URL with query parameters
	params := url.Values{}
	params.Add("apiKey", c.apiKey)
	params.Add("q", req.Query)

	if req.Sources != "" {
		params.Add("sources", req.Sources)
	}
	if req.Domains != "" {
		params.Add("domains", req.Domains)
	}
	if req.From != "" {
		params.Add("from", req.From)
	}
	if req.To != "" {
		params.Add("to", req.To)
	}
	if req.Language != "" {
		params.Add("language", req.Language)
	}
	if req.SortBy != "" {
		params.Add("sortBy", req.SortBy)
	}
	if req.PageSize > 0 {
		params.Add("pageSize", strconv.Itoa(req.PageSize))
	}
	if req.Page > 0 {
		params.Add("page", strconv.Itoa(req.Page))
	}

	requestURL := fmt.Sprintf("%s?%s", everythingURL, params.Encode())

	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("newsapi returned status: %d", resp.StatusCode)
	}

	var apiResp NewsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if apiResp.Status != "ok" {
		return nil, 0, fmt.Errorf("newsapi error: %s", apiResp.Status)
	}

	// Convert to domain models
	news := make([]*searchService.News, len(apiResp.Articles))
	for i, article := range apiResp.Articles {
		publishedAt, _ := time.Parse(time.RFC3339, article.PublishedAt)

		news[i] = &searchService.News{
			ID:          0, // Will be set by database
			Source:      article.Source.Name,
			Author:      article.Author,
			Title:       article.Title,
			Description: article.Description,
			URL:         article.URL,
			ImageURL:    article.URLToImage,
			PublishedAt: publishedAt,
			Content:     article.Content,
		}
	}

	return news, apiResp.TotalResults, nil
}

func (c *Client) GetTopHeadlines(ctx context.Context, req *searchService.TopHeadlinesRequest) ([]*searchService.News, int, error) {
	params := url.Values{}
	params.Add("apiKey", c.apiKey)

	if req.Country != "" {
		params.Add("country", req.Country)
	}
	if req.Category != "" {
		params.Add("category", req.Category)
	}
	if req.Sources != "" {
		params.Add("sources", req.Sources)
	}
	if req.Query != "" {
		params.Add("q", req.Query)
	}
	if req.PageSize > 0 {
		params.Add("pageSize", strconv.Itoa(req.PageSize))
	}
	if req.Page > 0 {
		params.Add("page", strconv.Itoa(req.Page))
	}

	requestURL := fmt.Sprintf("%s?%s", topHeadlinesURL, params.Encode())

	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("newsapi returned status: %d", resp.StatusCode)
	}

	var apiResp NewsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if apiResp.Status != "ok" {
		return nil, 0, fmt.Errorf("newsapi error: %s", apiResp.Status)
	}

	// Convert to domain models
	news := make([]*searchService.News, len(apiResp.Articles))
	for i, article := range apiResp.Articles {
		publishedAt, _ := time.Parse(time.RFC3339, article.PublishedAt)

		news[i] = &searchService.News{
			ID:          0,
			Source:      article.Source.Name,
			Author:      article.Author,
			Title:       article.Title,
			Description: article.Description,
			URL:         article.URL,
			ImageURL:    article.URLToImage,
			PublishedAt: publishedAt,
			Content:     article.Content,
		}
	}

	return news, apiResp.TotalResults, nil
}
