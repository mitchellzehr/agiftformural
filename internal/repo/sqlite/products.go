package sqlite

import (
	"context"
	"database/sql"
	"errors"

	muralerrors "mural/internal/errors"
	"mural/internal/model"
)

type productRow struct {
	ID        string  `db:"id"`
	Name      string  `db:"name"`
	Price     float64 `db:"price"`
	CreatedAt string  `db:"created_at"`
}

func (r *Repos) ListProducts(ctx context.Context) ([]model.Product, error) {
	var rows []productRow
	if err := r.db.SelectContext(ctx, &rows, `SELECT id, name, price, created_at FROM products ORDER BY name`); err != nil {
		return nil, err
	}
	out := make([]model.Product, len(rows))
	for i := range rows {
		out[i] = rowToProduct(rows[i])
	}
	return out, nil
}

func (r *Repos) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	var row productRow
	err := r.db.GetContext(ctx, &row, `SELECT id, name, price, created_at FROM products WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, muralerrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	p := rowToProduct(row)
	return &p, nil
}

func rowToProduct(row productRow) model.Product {
	return model.Product{
		ID:        row.ID,
		Name:      row.Name,
		Price:     row.Price,
		CreatedAt: parseTime(row.CreatedAt),
	}
}
