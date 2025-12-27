package api

import (
	"context"
	"gonews/protos/pb"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) CheckNewArticles(ctx context.Context, req *pb.CheckNewArticlesRequest) (*pb.CheckNewArticlesResponse, error) {
	log.Printf("CheckNewArticles called with keyword: %s, last_check_time: %s", req.Keyword, req.LastCheckTime)

	if req.Keyword == "" {
		return nil, status.Error(codes.InvalidArgument, "keyword is required")
	}

	// Парсим время последней проверки
	var lastCheckTime time.Time
	var err error

	if req.LastCheckTime != "" {
		lastCheckTime, err = time.Parse(time.RFC3339, req.LastCheckTime)
		if err != nil {
			log.Printf("Failed to parse last_check_time: %v, using default (24 hours ago)", err)
			lastCheckTime = time.Now().Add(-24 * time.Hour)
		}
	} else {
		lastCheckTime = time.Now().Add(-24 * time.Hour)
	}

	// Получаем все подписки на это ключ. слово
	subscriptions, err := s.notifyService.GetSubscriptionsByKeyword(ctx, req.Keyword)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get subscriptions: "+err.Error())
	}

	var allNewArticles []*pb.News
	userArticles := make(map[uint64][]*pb.News)

	// Для каждой подписки проверяем новые статьи
	for _, sub := range subscriptions {
		log.Printf("Checking subscription for user %d, keyword: %s", sub.UserID, req.Keyword)

		newArticles, err := s.notifyService.CheckNewArticlesForSubscription(ctx, sub.UserID, req.Keyword, lastCheckTime)
		if err != nil {
			log.Printf("Failed to check new articles for user %d: %v", sub.UserID, err)
			continue
		}

		if len(newArticles) > 0 {
			log.Printf("Found %d new articles for user %d, keyword: %s", len(newArticles), sub.UserID, req.Keyword)

			// Конвертируем в protobuf
			pbArticles := make([]*pb.News, len(newArticles))
			for i, article := range newArticles {
				var publishedAtStr string
				if !article.PublishedAt.IsZero() {
					publishedAtStr = article.PublishedAt.Format(time.RFC3339)
				}

				pbArticles[i] = &pb.News{
					Source:      article.Source,
					Author:      article.Author,
					Title:       article.Title,
					Description: article.Description,
					Url:         article.URL,
					ImageUrl:    article.URLToImage,
					PublishedAt: publishedAtStr,
				}
			}

			userArticles[sub.UserID] = pbArticles
			allNewArticles = append(allNewArticles, pbArticles...)
		}
	}

	log.Printf("CheckNewArticles completed. Found %d new articles for %d users", len(allNewArticles), len(userArticles))

	return &pb.CheckNewArticlesResponse{
		NewArticles: allNewArticles,
		UserStats:   s.createUserStats(userArticles),
	}, nil
}

func (s *GRPCServer) createUserStats(userArticles map[uint64][]*pb.News) []*pb.UserArticleStats {
	var stats []*pb.UserArticleStats
	for userID, articles := range userArticles {
		stats = append(stats, &pb.UserArticleStats{
			UserId:        userID,
			ArticlesCount: int32(len(articles)),
			Articles:      articles,
		})
	}
	return stats
}
