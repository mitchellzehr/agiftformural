package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	muralerrors "mural/internal/errors"
	"mural/internal/model"
)

type orderRow struct {
	ID          string  `db:"id"`
	Status      string  `db:"status"`
	TotalAmount float64 `db:"total_amount"`
	Currency    string  `db:"currency"`
	CreatedAt   string  `db:"created_at"`
}

type orderItemRow struct {
	ID        string  `db:"id"`
	OrderID   string  `db:"order_id"`
	ProductID string  `db:"product_id"`
	Quantity  int     `db:"quantity"`
	Price     float64 `db:"price"`
}

func (r *Repos) CreateOrderWithItems(ctx context.Context, o *model.Order, items []model.OrderItem) error {
	if o.CreatedAt.IsZero() {
		o.CreatedAt = time.Now().UTC()
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO orders (id, status, total_amount, currency, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		o.ID, string(o.Status), o.TotalAmount, o.Currency, formatTime(o.CreatedAt),
	)
	if err != nil {
		return err
	}

	for i := range items {
		it := &items[i]
		_, err = tx.ExecContext(ctx, `
			INSERT INTO order_items (id, order_id, product_id, quantity, price)
			VALUES (?, ?, ?, ?, ?)`,
			it.ID, it.OrderId, it.ProductId, it.Quantity, it.Price,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetOrderByID returns ErrNotFound when the order does not exist.
func (r *Repos) GetOrderByID(ctx context.Context, id string) (*model.Order, []model.OrderItem, error) {
	var row orderRow
	err := r.db.GetContext(ctx, &row, `SELECT id, status, total_amount, currency, created_at FROM orders WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, muralerrors.ErrNotFound
	}
	if err != nil {
		return nil, nil, err
	}
	o := rowToOrder(row)
	items, err := r.orderItemsForOrderIDs(ctx, []string{id})
	if err != nil {
		return nil, nil, err
	}
	return &o, items, nil
}

func (r *Repos) ListOrders(ctx context.Context) ([]model.Order, []model.OrderItem, error) {
	var rows []orderRow
	if err := r.db.SelectContext(ctx, &rows, `
		SELECT id, status, total_amount, currency, created_at
		FROM orders ORDER BY created_at DESC`); err != nil {
		return nil, nil, err
	}
	orders := make([]model.Order, len(rows))
	ids := make([]string, len(rows))
	for i := range rows {
		orders[i] = rowToOrder(rows[i])
		ids[i] = rows[i].ID
	}
	items, err := r.orderItemsForOrderIDs(ctx, ids)
	if err != nil {
		return nil, nil, err
	}
	return orders, items, nil
}

func (r *Repos) orderItemsForOrderIDs(ctx context.Context, orderIDs []string) ([]model.OrderItem, error) {
	if len(orderIDs) == 0 {
		return []model.OrderItem{}, nil
	}
	query, args, err := sqlx.In(`
		SELECT id, order_id, product_id, quantity, price
		FROM order_items WHERE order_id IN (?) ORDER BY order_id, id`, orderIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	var itemRows []orderItemRow
	if err := r.db.SelectContext(ctx, &itemRows, query, args...); err != nil {
		return nil, err
	}
	out := make([]model.OrderItem, len(itemRows))
	for i := range itemRows {
		out[i] = rowToOrderItem(itemRows[i])
	}
	return out, nil
}

func rowToOrder(row orderRow) model.Order {
	return model.Order{
		ID:          row.ID,
		Status:      model.OrderStatus(row.Status),
		TotalAmount: row.TotalAmount,
		Currency:    row.Currency,
		CreatedAt:   parseTime(row.CreatedAt),
	}
}

func rowToOrderItem(row orderItemRow) model.OrderItem {
	return model.OrderItem{
		ID:        row.ID,
		OrderId:   row.OrderID,
		ProductId: row.ProductID,
		Quantity:  row.Quantity,
		Price:     row.Price,
	}
}
