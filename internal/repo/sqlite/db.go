package sqlite

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// Open returns a sqlx DB using the modernc CGO-free SQLite driver.
func Open(path string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)", path)
	db, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

// InitSchema applies DDL (CREATE TABLE IF NOT EXISTS).
func InitSchema(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, schema)
	return err
}
