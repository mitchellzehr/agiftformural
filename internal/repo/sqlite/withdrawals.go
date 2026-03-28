package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	muralerrors "mural/internal/errors"
	"mural/internal/model"
)

type withdrawalRow struct {
	ID              string  `db:"id"`
	OrderID         string  `db:"order_id"`
	MuralTransferID string  `db:"mural_transfer_id"`
	Amount          float64 `db:"amount"`
	SourceCurrency  string  `db:"source_currency"`
	DestCurrency    string  `db:"dest_currency"`
	Status          string  `db:"status"`
	CreatedAt       string  `db:"created_at"`
}

func (r *Repos) ListWithdrawals(ctx context.Context) ([]model.Withdrawal, error) {
	var rows []withdrawalRow
	if err := r.db.SelectContext(ctx, &rows, `
		SELECT id, order_id, mural_transfer_id, amount, source_currency, dest_currency, status, created_at
		FROM withdrawals ORDER BY created_at DESC`); err != nil {
		return nil, err
	}
	out := make([]model.Withdrawal, len(rows))
	for i := range rows {
		out[i] = rowToWithdrawal(rows[i])
	}
	return out, nil
}

func (r *Repos) GetWithdrawal(ctx context.Context, id string) (*model.Withdrawal, error) {
	var row withdrawalRow
	err := r.db.GetContext(ctx, &row, `
		SELECT id, order_id, mural_transfer_id, amount, source_currency, dest_currency, status, created_at
		FROM withdrawals WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, muralerrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	w := rowToWithdrawal(row)
	return &w, nil
}

func (r *Repos) CreateWithdrawal(ctx context.Context, w *model.Withdrawal) error {
	if w.CreatedAt.IsZero() {
		w.CreatedAt = time.Now().UTC()
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO withdrawals (id, order_id, mural_transfer_id, amount, source_currency, dest_currency, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		w.ID, w.OrderID, w.MuralTransferID, w.Amount, w.SourceCurrency, w.DestCurrency, string(w.Status), formatTime(w.CreatedAt),
	)
	return err
}

func rowToWithdrawal(row withdrawalRow) model.Withdrawal {
	return model.Withdrawal{
		ID:              row.ID,
		OrderID:         row.OrderID,
		MuralTransferID: row.MuralTransferID,
		Amount:          row.Amount,
		SourceCurrency:  row.SourceCurrency,
		DestCurrency:    row.DestCurrency,
		Status:          model.WithdrawalStatus(row.Status),
		CreatedAt:       parseTime(row.CreatedAt),
	}
}
