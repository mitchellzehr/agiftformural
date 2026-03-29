package sqlite

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

// defaultProductSeeds populate an empty catalog on first run; ids are stable for docs and curl examples.
var defaultProductSeeds = []struct {
	ID    string
	Name  string
	Price float64
}{
	{"prod-poster-1", "Community poster", 29.99},
	{"prod-stickers-1", "Sticker sheet set", 12.50},
	{"prod-pin-1", "Enamel pin", 8.00},
}

// SeedDefaultProducts inserts default catalog rows when missing (INSERT OR IGNORE by primary key).
// Safe to call on every process start.
func SeedDefaultProducts(ctx context.Context, db *sqlx.DB) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	now := time.Now().UTC().Format(time.RFC3339)
	const q = `INSERT OR IGNORE INTO products (id, name, price, created_at) VALUES (?, ?, ?, ?)`
	for _, p := range defaultProductSeeds {
		if _, err := tx.ExecContext(ctx, q, p.ID, p.Name, p.Price, now); err != nil {
			return err
		}
	}
	return tx.Commit()
}
