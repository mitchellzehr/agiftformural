package handler

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	muralerrors "mural/internal/errors"
	"mural/internal/model"
)

type paymentWebhookRequest struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

func (s *Server) PaymentWebhook(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var req paymentWebhookRequest
	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}
	orderID := strings.TrimSpace(req.OrderID)
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "order_id required"})
		return
	}
	o, _, err := s.orderSvc.GetOrderByID(r.Context(), orderID)
	if err != nil {
		if errors.Is(err, muralerrors.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "order_not_found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	// FLOAT VALUE COMPARISON
	const eps = 1e-6
	if math.Abs(req.Amount-o.TotalAmount) > eps {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "amount does not match order total"})
		return
	}

	p := &model.Payment{
		ID:        uuid.NewString(),
		OrderID:   orderID,
		Amount:    req.Amount,
		Currency:  o.Currency,
		Status:    model.PaymentStatusReceived,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.orderSvc.RecordPayment(r.Context(), p); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
