package api

import (
	"context"
	"gonews/notify_service/internal/models"
	"gonews/protos/pb"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	log.Printf("SendNotification called for user %d", req.UserId)

	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	if req.Message == "" && len(req.Articles) == 0 {
		return nil, status.Error(codes.InvalidArgument, "either message or articles must be provided")
	}

	// Если есть статьи, отправляем уведомления для каждой
	if len(req.Articles) > 0 {
		log.Printf("Sending notifications for %d articles to user %d", len(req.Articles), req.UserId)

		// Определяем тему уведомления из первой статьи или запроса
		notificationTopic := "news"
		if req.Message != "" {
			notificationTopic = req.Message
		} else if len(req.Articles) > 0 && req.Articles[0].Title != "" {
			// Берем ключевые слова из заголовка
			notificationTopic = "new_articles"
		}

		// Отправляем каждую статью
		for _, pbArticle := range req.Articles {
			// Парсим время публикации
			var publishedAt time.Time
			if pbArticle.PublishedAt != "" {
				publishedAt, _ = time.Parse(time.RFC3339, pbArticle.PublishedAt)
			}

			// Создаем модель новости
			article := models.News{
				Source:      pbArticle.Source,
				Author:      pbArticle.Author,
				Title:       pbArticle.Title,
				Description: pbArticle.Description,
				URL:         pbArticle.Url,
				URLToImage:  pbArticle.ImageUrl,
				PublishedAt: publishedAt,
			}

			// Отправляем уведомление
			err := s.notifyService.SendNotification(ctx, req.UserId, notificationTopic, article)
			if err != nil {
				log.Printf("Failed to send notification for article '%s': %v", pbArticle.Title, err)
				// Продолжаем отправлять остальные уведомления
				continue
			}

			log.Printf("Notification sent for article: %s", pbArticle.Title)
		}

		return &pb.SendNotificationResponse{
			Success:   true,
			SentCount: int32(len(req.Articles)),
			Message:   "Notifications sent successfully",
		}, nil
	}

	// Если есть только сообщение без статей
	if req.Message != "" {
		log.Printf("Sending simple message to user %d: %s", req.UserId, req.Message)

		// Создаем фиктивную статью для структуры уведомления
		article := models.News{
			Source:      "System",
			Author:      "Notification Service",
			Title:       req.Message,
			Description: req.Message,
			URL:         "",
			URLToImage:  "",
			PublishedAt: time.Now(),
		}

		err := s.notifyService.SendNotification(ctx, req.UserId, "system_notification", article)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to send notification: "+err.Error())
		}

		return &pb.SendNotificationResponse{
			Success:   true,
			SentCount: 1,
			Message:   "Message notification sent successfully",
		}, nil
	}

	return &pb.SendNotificationResponse{
		Success: false,
		Message: "No content to send",
	}, nil
}
