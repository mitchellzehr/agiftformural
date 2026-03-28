package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	muralerrors "mural/internal/errors"
	"mural/internal/model"
)

type paymentRow struct {
	ID        string  `db:"id"`
	OrderID   string  `db:"order_id"`
	Amount    float64 `db:"amount"`
	Currency  string  `db:"currency"`
	Status    string  `db:"status"`
	CreatedAt string  `db:"created_at"`
}

func (r *Repos) CreatePayment(ctx context.Context, p *model.Payment) error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO payments (id, order_id, amount, currency, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		p.ID, p.OrderID, p.Amount, p.Currency, string(p.Status), formatTime(p.CreatedAt),
	)
	return err
}

func (r *Repos) GetPaymentByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	var row paymentRow
	err := r.db.GetContext(ctx, &row, `
		SELECT id, order_id, amount, currency, status, created_at
		FROM payments WHERE order_id = ? ORDER BY created_at DESC LIMIT 1`, orderID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, muralerrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	p := rowToPayment(row)
	return &p, nil
}

func rowToPayment(row paymentRow) model.Payment {
	return model.Payment{
		ID:        row.ID,
		OrderID:   row.OrderID,
		Amount:    row.Amount,
		Currency:  row.Currency,
		Status:    model.PaymentStatus(row.Status),
		CreatedAt: parseTime(row.CreatedAt),
	}
}
