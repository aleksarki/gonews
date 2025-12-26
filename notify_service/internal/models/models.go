package models

import "time"

type News struct {
	Source      string    `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
}

type NotificationMessage struct {
	EventID    string    `json:"event_id"`
	EventType  string    `json:"event_type"`
	NotifTopic string    `json:"notif_topic"`
	Article    News      `json:"article"`
	Timestamp  time.Time `json:"timestamp"`
}

type Subscription struct {
	ID      uint64 `json:"id"`
	UserID  uint64 `json:"user_id"`
	Keyword string `json:"keyword"`
}
