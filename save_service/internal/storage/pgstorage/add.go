package pgstorage

import (
	"context"
	"gonews/save_service/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func (storage *PGStorage) UpsertNews(ctx context.Context, news []*models.News) error {
	news_ := lo.Map(news, func(n *models.News, _ int) *News {
		return &News{
			Source:      n.Source,
			Author:      n.Author,
			Title:       n.Title,
			Description: n.Description,
			URL:         n.URL,
			ImageURL:    n.ImageURL,
			PublishedAt: n.PublishedAt,
		}
	})
	query := squirrel.Insert("news").
		Columns("source", "author", "title", "description", "url", "image_url", "published_at").
		PlaceholderFormat(squirrel.Dollar)
	for _, n := range news_ {
		query = query.Values(n.Source, n.Author, n.Title, n.Description, n.URL, n.ImageURL, n.PublishedAt)
	}
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
