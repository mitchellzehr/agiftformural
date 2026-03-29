package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	muralerrors "mural/internal/errors"
	"mural/internal/model"
)

// OrderSvc is what order and payment webhook handlers need from the application layer.
type OrderSvc interface {
	ListOrders(ctx context.Context) ([]model.Order, []model.OrderItem, error)
	GetOrderByID(ctx context.Context, id string) (*model.Order, []model.OrderItem, error)
	CreateOrder(ctx context.Context, order *model.Order, items []model.OrderItem) error
	RecordPayment(ctx context.Context, p *model.Payment) error
}

type orderResponse struct {
	ID          string  `json:"id"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_amount"`
	Currency    string  `json:"currency"`
	CreatedAt   string  `json:"created_at"`
}

type orderItemResponse struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type createOrderRequest struct {
	ID       string                 `json:"id"`
	Currency string                 `json:"currency"`
	Items    []createOrderLineInput `json:"items"`
}

type createOrderLineInput struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (s *Server) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, lineItems, err := s.orderSvc.ListOrders(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	byOrder := groupOrderItems(lineItems)
	out := make([]map[string]any, 0, len(orders))
	for _, o := range orders {
		out = append(out, map[string]any{
			"order": orderFromModel(o),
			"items": orderItemsToResponse(byOrder[o.ID]),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"orders": out})
}

func (s *Server) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "missing order id"})
		return
	}
	o, items, err := s.orderSvc.GetOrderByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, muralerrors.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not_found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"order": orderFromModel(*o),
		"items": orderItemsToResponse(items),
	})
}

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var req createOrderRequest

	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}
	req.Currency = strings.TrimSpace(req.Currency)
	if req.Currency == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "currency required"})
		return
	}
	orderID := strings.TrimSpace(req.ID)
	if orderID == "" {
		orderID = uuid.NewString()
	}

	now := time.Now().UTC()
	items := make([]model.OrderItem, len(req.Items))
	for i := range req.Items {
		line := req.Items[i]
		pid := strings.TrimSpace(line.ProductID)
		if pid == "" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": fmt.Sprintf("items[%d]: product_id required", i)})
			return
		}
		if line.Quantity <= 0 {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": fmt.Sprintf("items[%d]: quantity must be positive", i)})
			return
		}
		lid := strings.TrimSpace(line.ID)
		if lid == "" {
			lid = uuid.NewString()
		}
		items[i] = model.OrderItem{
			ID:        lid,
			ProductId: pid,
			Quantity:  line.Quantity,
		}
	}

	order := &model.Order{
		ID:        orderID,
		Status:    model.OrderStatusPendingPayment,
		Currency:  req.Currency,
		CreatedAt: now,
	}
	if err := s.orderSvc.CreateOrder(r.Context(), order, items); err != nil {
		if errors.Is(err, muralerrors.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "product_not_found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	o, resolved, err := s.orderSvc.GetOrderByID(r.Context(), orderID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"order": orderFromModel(*o),
		"items": orderItemsToResponse(resolved),
	})
}

// TRANSLATION HELPERS

func orderFromModel(o model.Order) orderResponse {
	return orderResponse{
		ID:          o.ID,
		Status:      string(o.Status),
		TotalAmount: o.TotalAmount,
		Currency:    o.Currency,
		CreatedAt:   o.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func groupOrderItems(items []model.OrderItem) map[string][]model.OrderItem {
	m := make(map[string][]model.OrderItem)
	for _, it := range items {
		m[it.OrderId] = append(m[it.OrderId], it)
	}
	return m
}

func orderItemsToResponse(items []model.OrderItem) []orderItemResponse {
	out := make([]orderItemResponse, 0, len(items))
	for _, it := range items {
		out = append(out, orderItemResponse{
			ID:        it.ID,
			OrderID:   it.OrderId,
			ProductID: it.ProductId,
			Quantity:  it.Quantity,
			Price:     it.Price,
		})
	}
	return out
}
