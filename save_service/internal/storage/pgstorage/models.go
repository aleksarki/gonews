package pgstorage

import "time"

type User struct {
	ID   uint64 `db:"id"`
	Name string `db:"name"`
}

type News struct {
	ID          uint64    `db:"id"`
	Source      string    `db:"source"`
	Author      string    `db:"author"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	URL         string    `db:"url"`
	ImageURL    string    `db:"image_url"`
	PublishedAt time.Time `db:"published_at"`
}

type UserToFavouriteNews struct {
	ID     uint64 `db:"id"`
	UserID uint64 `db:"user_id"`
	NewsID uint64 `db:"news_id"`
}

type UserToSeenNews struct {
	ID     uint64    `db:"id"`
	UserID uint64    `db:"user_id"`
	NewsID uint64    `db:"news_id"`
	SeenAt time.Time `db:"seen_at"`
}

type Subscription struct {
	ID      uint64 `db:"id"`
	UserID  uint64 `db:"user_id"`
	Keyword string `db:"keyword"`
}
