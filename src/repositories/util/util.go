package util

import (
	"context"
	"embed"
	"io/fs"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

func LoadSQLQueries(sqlFiles embed.FS) (map[string]string, error) {
	queries := make(map[string]string)

	err := fs.WalkDir(sqlFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		queryName := filepath.Base(path[:len(path)-len(filepath.Ext(path))])
		content, err := fs.ReadFile(sqlFiles, path)
		if err != nil {
			return err
		}

		queries[queryName] = string(content)
		return nil
	})

	return queries, err
}

type DBQueryParams struct {
	DBConnection *sqlx.DB
	SqlQuery     string
	Variables    map[string]any
}

func DBQuery[T any](ctx context.Context, p DBQueryParams) ([]T, error) {
	query, args, err := sqlx.Named(p.SqlQuery, p.Variables)
	if err != nil {
		return nil, err
	}
	var results []T
	err = p.DBConnection.SelectContext(ctx, &results, p.DBConnection.Rebind(query), args...)
	if err != nil {
		return nil, err
	}
	return results, err
}
