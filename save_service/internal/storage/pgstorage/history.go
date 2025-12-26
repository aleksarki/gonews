package pgstorage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

// AddToSearchHistory - добавляем запись в историю поиска
func (storage *PGStorage) AddToSearchHistory(ctx context.Context, userID uint64, query string, results []uint64) error {
	// Сохраняем запрос
	searchQuery := squirrel.Insert("search_history").
		Columns("user_id", "query", "searched_at").
		Values(userID, query, time.Now()).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := searchQuery.ToSql()
	if err != nil {
		return errors.Wrap(err, "query generation error")
	}

	var searchID uint64
	err = storage.DB.QueryRow(ctx, queryText, args...).Scan(&searchID)
	if err != nil {
		return errors.Wrap(err, "query execution error")
	}

	// Сохраняем результаты поиска
	if len(results) > 0 {
		resultsJSON, err := json.Marshal(results)
		if err != nil {
			return errors.Wrap(err, "failed to marshal results")
		}

		resultsQuery := squirrel.Insert("search_results").
			Columns("search_id", "news_ids").
			Values(searchID, resultsJSON).
			PlaceholderFormat(squirrel.Dollar)

		resultsQueryText, resultsArgs, err := resultsQuery.ToSql()
		if err != nil {
			return errors.Wrap(err, "results query generation error")
		}

		_, err = storage.DB.Exec(ctx, resultsQueryText, resultsArgs...)
		if err != nil {
			return errors.Wrap(err, "results query execution error")
		}
	}

	return nil
}

// GetSearchHistory - получаем историю поиска пользователя
func (storage *PGStorage) GetSearchHistory(ctx context.Context, userID uint64) ([]string, error) {
	query := squirrel.Select("query").
		From("search_history").
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("searched_at DESC").
		Limit(50).
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

	var queries []string
	for rows.Next() {
		var query string
		err := rows.Scan(&query)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		queries = append(queries, query)
	}

	return queries, nil
}

// MarkNewsAsSeen - отмечаем новость как просмотренную
func (storage *PGStorage) MarkNewsAsSeen(ctx context.Context, userID, newsID uint64) error {
	query := squirrel.Insert("user_to_seen_news").
		Columns("user_id", "news_id").
		Values(userID, newsID).
		Suffix("ON CONFLICT (user_id, news_id) DO UPDATE SET seen_at = CURRENT_TIMESTAMP").
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
