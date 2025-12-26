package pgstorage

import (
	"context"
	"gonews/save_service/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func (storage *PGStorage) AddFavourite(ctx context.Context, userID, newsID uint64) error {
	// Проверяем существование новости
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM news WHERE id = $1)`
	err := storage.DB.QueryRow(ctx, checkQuery, newsID).Scan(&exists)
	if err != nil || !exists {
		return errors.New("news not found")
	}

	query := squirrel.Insert("user_to_favourite_news").
		Columns("user_id", "news_id").
		Values(userID, newsID).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "query generation error")
	}

	_, err = storage.DB.Exec(ctx, queryText, args...)
	if err != nil {
		return errors.Wrap(err, "query execution error")
	}

	return nil
}

// GetFavourites - получаем избранные новости пользователя
func (storage *PGStorage) GetFavourites(ctx context.Context, userID uint64) ([]*models.News, error) {
	query := squirrel.Select("n.id", "n.source", "n.author", "n.title", "n.description",
		"n.url", "n.image_url", "n.published_at").
		From("user_to_favourite_news f").
		Join("news n ON f.news_id = n.id").
		Where(squirrel.Eq{"f.user_id": userID}).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query generation error")
	}

	rows, err := storage.DB.Query(ctx, queryText, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	defer rows.Close()

	var news []*models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(&n.ID, &n.Source, &n.Author, &n.Title, &n.Description,
			&n.URL, &n.ImageURL, &n.PublishedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		news = append(news, &n)
	}

	return news, nil
}
