package models

import "time"

type User struct {
	ID   uint64
	Name string
}

type News struct {
	ID          uint64
	Source      string
	Author      string
	Title       string
	Description string
	URL         string
	ImageURL    string
	PublishedAt time.Time
}

type UserToFavouriteNews struct {
	ID     uint64
	UserID uint64
	NewsID uint64
}

type UserToSeenNews struct {
	ID     uint64
	UserID uint64
	NewsID uint64
	SeenAt time.Time
}

type Subscription struct {
	ID      uint64
	UserID  uint64
	Keyword string
}
