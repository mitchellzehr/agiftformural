package service

import (
	"context"

	"mural/internal/model"
)

// WithdrawalService exposes payout / withdrawal reads and creates.
type WithdrawalService struct {
	withdrawals WithdrawalStore
}

// NewWithdrawalService builds a withdrawal service over a WithdrawalStore (e.g. repos.Withdrawals).
func NewWithdrawalService(withdrawals WithdrawalStore) *WithdrawalService {
	return &WithdrawalService{withdrawals: withdrawals}
}

// ListWithdrawals returns all withdrawals (ordering defined by the store).
func (s *WithdrawalService) ListWithdrawals(ctx context.Context) ([]model.Withdrawal, error) {
	return s.withdrawals.ListWithdrawals(ctx)
}

// GetWithdrawal returns one withdrawal or an error (including ErrNotFound from the store).
func (s *WithdrawalService) GetWithdrawal(ctx context.Context, id string) (*model.Withdrawal, error) {
	return s.withdrawals.GetWithdrawal(ctx, id)
}

// CreateWithdrawal persists a new withdrawal row (append-only).
func (s *WithdrawalService) CreateWithdrawal(ctx context.Context, w *model.Withdrawal) error {
	return s.withdrawals.CreateWithdrawal(ctx, w)
}
