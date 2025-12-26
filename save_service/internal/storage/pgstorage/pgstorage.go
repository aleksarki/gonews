package pgstorage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PGStorage struct {
	DB *pgxpool.Pool
}

func (storage *PGStorage) InitTables() error {
	sql := `
		CREATE TABLE IF NOT EXISTS Users(
			id    SERIAL        PRIMARY KEY,
			name  VARCHAR(255)  NOT NULL
		);

		CREATE TABLE IF NOT EXISTS News (
			id            SERIAL         PRIMARY KEY,
			source        VARCHAR(255),
			author        VARCHAR(255),
			title         VARCHAR(500)   NOT NULL,
			description   TEXT           NOT NULL,
			url           VARCHAR(2048)  NOT NULL,
			image_url     VARCHAR(2048),
			published_at  TIMESTAMP,
			UNIQUE(url)
		);

		CREATE TABLE IF NOT EXISTS UserToFavouriteNews (
			id       SERIAL  PRIMARY KEY,
			user_id  INT     NOT NULL,
			news_id  INT     NOT NULL,

			CONSTRAINT fk_user_favourite_news_user
				FOREIGN KEY (user_id)
				REFERENCES Users(id)
				ON DELETE CASCADE,

			CONSTRAINT fk_user_favourite_news_news
				FOREIGN KEY (news_id)
				REFERENCES News(id)
				ON DELETE CASCADE,

			CONSTRAINT unique_user_favourite_news_user_news 
				UNIQUE (user_id, news_id)
		);

		CREATE TABLE IF NOT EXISTS UserToSeenNews (
			id       SERIAL     PRIMARY KEY,
			user_id  INT        NOT NULL,
			news_id  INT        NOT NULL,
			seen_at  TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,

			CONSTRAINT fk_user_seen_news_user
				FOREIGN KEY (user_id)
				REFERENCES Users(id)
				ON DELETE CASCADE,

			CONSTRAINT fk_user_seen_news_news
				FOREIGN KEY (news_id)
				REFERENCES News(id)
				ON DELETE CASCADE,

			CONSTRAINT unique_user_seen_news_user_news
				UNIQUE (user_id, news_id)
		);

		CREATE TABLE IF NOT EXISTS user_subscriptions (
			id      SERIAL       PRIMARY KEY,
			user_id BIGINT       NOT NULL,
			keyword VARCHAR(100) NOT NULL,
			
			CONSTRAINT fk_subscriptions_user 
				FOREIGN KEY (user_id) 
				REFERENCES Users(id) 
				ON DELETE CASCADE,

			CONSTRAINT unique_user_subscriptions_user_keyword 
				UNIQUE (user_id, keyword)
		);

		CREATE TABLE IF NOT EXISTS search_history (
			id          SERIAL      PRIMARY KEY,
			user_id     BIGINT      NOT NULL,
			query       TEXT        NOT NULL,
			searched_at TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
			
			CONSTRAINT fk_search_history_user 
				FOREIGN KEY (user_id) 
				REFERENCES Users(id) 
				ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS search_results (
			id         SERIAL  PRIMARY KEY,
			search_id  BIGINT  NOT NULL,
			news_ids   JSONB   NOT NULL,
			
			CONSTRAINT fk_search_results_search 
				FOREIGN KEY (search_id) 
				REFERENCES search_history(id) 
				ON DELETE CASCADE
		);
	`
	_, err := storage.DB.Exec(context.Background(), sql)
	if err != nil {
		return errors.Wrap(err, "table initialization error")
	}

	return nil
}

func NewPgstorage(connectionString string) (*PGStorage, error) {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "config parsing error")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, "connection error")
	}

	storage := &PGStorage{DB: db}
	err = storage.InitTables()
	if err != nil {
		return nil, err
	}

	return storage, nil
}
