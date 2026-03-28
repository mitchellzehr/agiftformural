package service

import (
	"context"
	"fmt"

	"mural/internal/model"
)

// OrderService coordinates orders, catalog-backed line pricing, and payment rows.
type OrderService struct {
	orders   OrderStore
	products ProductReader
	payments PaymentStore
}

// NewOrderService wires order, catalog, and payment persistence.
func NewOrderService(orders OrderStore, products ProductReader, payments PaymentStore) *OrderService {
	return &OrderService{orders: orders, products: products, payments: payments}
}

// ListOrders returns all orders and a flat slice of all line items (same contract as the store).
func (s *OrderService) ListOrders(ctx context.Context) ([]model.Order, []model.OrderItem, error) {
	return s.orders.ListOrders(ctx)
}

// GetOrderByID returns the order and its line items, or ErrNotFound from the store.
func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*model.Order, []model.OrderItem, error) {
	return s.orders.GetOrderByID(ctx, id)
}

// CreateOrder validates line items against the catalog, snapshots unit Price from each product,
// sets OrderId on each line, recomputes TotalAmount from Σ(price × quantity), and persists.
func (s *OrderService) CreateOrder(ctx context.Context, order *model.Order, items []model.OrderItem) error {
	if order == nil {
		return fmt.Errorf("service: nil order")
	}
	resolved := make([]model.OrderItem, len(items))
	var sum float64
	for i := range items {
		it := items[i]
		if it.ProductId == "" {
			return fmt.Errorf("service: line %d: empty product id", i)
		}
		if it.Quantity <= 0 {
			return fmt.Errorf("service: line %d: quantity must be positive", i)
		}
		p, err := s.products.GetProduct(ctx, it.ProductId)
		if err != nil {
			return err
		}
		it.OrderId = order.ID
		it.Price = p.Price
		sum += it.Price * float64(it.Quantity)
		resolved[i] = it
	}
	order.TotalAmount = sum
	return s.orders.CreateOrderWithItems(ctx, order, resolved)
}

// RecordPayment appends a payment row (e.g. after a webhook confirms funds).
func (s *OrderService) RecordPayment(ctx context.Context, p *model.Payment) error {
	if p == nil {
		return fmt.Errorf("service: nil payment")
	}
	return s.payments.CreatePayment(ctx, p)
}

// GetPaymentForOrder returns the latest payment row for the order, or ErrNotFound from the store.
func (s *OrderService) GetPaymentForOrder(ctx context.Context, orderID string) (*model.Payment, error) {
	return s.payments.GetPaymentByOrderID(ctx, orderID)
}
