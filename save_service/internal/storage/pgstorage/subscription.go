package pgstorage

import (
	"context"
	"gonews/save_service/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

// Subscribe - добавляем подписку пользователя
func (storage *PGStorage) Subscribe(ctx context.Context, userID uint64, keyword string) error {
	query := squirrel.Insert("user_subscriptions").
		Columns("user_id", "keyword").
		Values(userID, keyword).
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

// GetSubscriptions - получаем все подписки
func (storage *PGStorage) GetSubscriptions(ctx context.Context) ([]*models.Subscription, error) {
	query := squirrel.Select("id", "user_id", "keyword").
		From("user_subscriptions").
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

	var subscriptions []*models.Subscription
	for rows.Next() {
		var s models.Subscription
		err := rows.Scan(&s.ID, &s.UserID, &s.Keyword)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		subscriptions = append(subscriptions, &s)
	}

	return subscriptions, nil
}
