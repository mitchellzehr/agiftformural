package sqlite

import (
	"github.com/jmoiron/sqlx"

	"mural/internal/service"
)

type Repos struct {
	db *sqlx.DB
}

// NewRepos builds service.Repos backed by SQLite.
func NewRepos(db *sqlx.DB) service.Repos {
	r := &Repos{db: db}
	return service.Repos{
		Products:    r,
		Orders:      r,
		Payments:    r,
		Withdrawals: r,
	}
}
