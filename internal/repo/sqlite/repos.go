package sqlite

import (
	"github.com/jmoiron/sqlx"
)

type Repos struct {
	db *sqlx.DB
}

// NewRepos builds service.Repos backed by SQLite.
func NewRepos(db *sqlx.DB) Repos {
	return Repos{db: db}
}
