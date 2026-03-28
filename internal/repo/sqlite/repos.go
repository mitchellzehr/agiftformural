package sqlite

import (
	"github.com/jmoiron/sqlx"

	"mural/internal/repo"
)

type Repos struct {
	db *sqlx.DB
}

// NewRepos builds repo.Repos backed by SQLite.
func NewRepos(db *sqlx.DB) repo.Repos {
	r := &Repos{db: db}
	return repo.Repos{
		Products:    r,
		Orders:      r,
		Payments:    r,
		Withdrawals: r,
	}
}
