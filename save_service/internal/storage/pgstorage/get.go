package pgstorage

import (
	"context"
	"gonews/save_service/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func (storage *PGStorage) GetNewsByIDs(ctx context.Context, IDs []uint64) ([]*models.News, error) {
	query := squirrel.Select("id", "source", "author", "title", "description", "url", "image_url", "published_at").
		From("News").
		Where(squirrel.Eq{"id": IDs}).
		PlaceholderFormat(squirrel.Dollar)
	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query generation error")
	}

	rows, err := storage.DB.Query(ctx, queryText, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}

	var news []*models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(&n.ID, &n.Source, &n.Author, &n.Title, &n.Description, &n.URL, &n.ImageURL, &n.PublishedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		news = append(news, &n)
	}
	return news, nil
}
