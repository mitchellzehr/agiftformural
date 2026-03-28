package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	muralerrors "mural/internal/errors"
	"mural/internal/model"
)

// WithdrawalReader is what withdrawal list/get handlers need from the store.
type WithdrawalReader interface {
	ListWithdrawals(ctx context.Context) ([]model.Withdrawal, error)
	GetWithdrawal(ctx context.Context, id string) (*model.Withdrawal, error)
}

type withdrawalResponse struct {
	ID              string  `json:"id"`
	OrderID         string  `json:"order_id"`
	MuralTransferID string  `json:"mural_transfer_id"`
	Amount          float64 `json:"amount"`
	SourceCurrency  string  `json:"source_currency"`
	DestCurrency    string  `json:"dest_currency"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
}

func (s *Server) ListWithdrawals(w http.ResponseWriter, r *http.Request) {
	items, err := s.withdrawals.ListWithdrawals(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	out := make([]withdrawalResponse, 0, len(items))
	for _, wdr := range items {
		out = append(out, withdrawalFromModel(wdr))
	}
	writeJSON(w, http.StatusOK, map[string]any{"withdrawals": out})
}

func (s *Server) GetWithdrawal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "missing withdrawal id"})
		return
	}
	wdr, err := s.withdrawals.GetWithdrawal(r.Context(), id)
	if err != nil {
		if errors.Is(err, muralerrors.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not_found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"withdrawal": withdrawalFromModel(*wdr)})
}

func withdrawalFromModel(w model.Withdrawal) withdrawalResponse {
	return withdrawalResponse{
		ID:              w.ID,
		OrderID:         w.OrderID,
		MuralTransferID: w.MuralTransferID,
		Amount:          w.Amount,
		SourceCurrency:  w.SourceCurrency,
		DestCurrency:    w.DestCurrency,
		Status:          string(w.Status),
		CreatedAt:       w.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}
