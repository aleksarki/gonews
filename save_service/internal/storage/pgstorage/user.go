package pgstorage

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

// CreateUser - создаем пользователя
func (storage *PGStorage) CreateUser(ctx context.Context, name string) (uint64, error) {
	query := squirrel.Insert("users").
		Columns("name").
		Values(name).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "query generation error")
	}

	var id uint64
	err = storage.DB.QueryRow(ctx, queryText, args...).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "query execution error")
	}

	return id, nil
}
