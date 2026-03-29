package service

import (
	"context"
	"fmt"

	"mural/internal/model"

	"github.com/google/uuid"
)

// MuralClient is implemented by mural.Client (HTTP) and mural.StubClient (no network).
type MuralClient interface {
	CreateTransfer(ctx context.Context, amount float64) (transferID string, err error)
}

// OrderService coordinates orders, catalog-backed line pricing, payments, and Mural payout stubs.
type OrderService struct {
	orders      OrderStore
	products    ProductReader
	payments    PaymentStore
	withdrawals WithdrawalStore
	muralClient MuralClient
}

// NewOrderService wires order, catalog, payment, withdrawal persistence, and the Mural client (or stub).
func NewOrderService(
	orders OrderStore,
	products ProductReader,
	payments PaymentStore,
	withdrawals WithdrawalStore,
	muralClient MuralClient,
) *OrderService {
	return &OrderService{
		orders:      orders,
		products:    products,
		payments:    payments,
		withdrawals: withdrawals,
		muralClient: muralClient,
	}
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

// RecordPayment persists the payment, marks the order paid, and calls Mural CreateTransfer (stub or HTTP client).
// If the order is already paid, returns nil without inserting another payment (idempotent for duplicate webhooks).
func (s *OrderService) RecordPayment(ctx context.Context, p *model.Payment) error {
	if p == nil {
		return fmt.Errorf("service: nil payment")
	}
	order, _, err := s.orders.GetOrderByID(ctx, p.OrderID)
	if err != nil {
		return err
	}
	if order.Status == model.OrderStatusPaid {
		return nil
	}
	if err := s.payments.CreatePayment(ctx, p); err != nil {
		return err
	}
	if err := s.orders.UpdateOrderStatus(ctx, p.OrderID, model.OrderStatusPaid); err != nil {
		return err
	}
	transferID, err := s.muralClient.CreateTransfer(ctx, p.Amount)
	if err != nil {
		return err
	}
	return s.withdrawals.CreateWithdrawal(ctx, &model.Withdrawal{
		ID:              uuid.NewString(),
		OrderID:         p.OrderID,
		MuralTransferID: transferID,
		Amount:          p.Amount,
		SourceCurrency:  order.Currency,
		DestCurrency:    "COP",
		Status:          model.WithdrawalStatusPending,
	})
}

// GetPaymentForOrder returns the latest payment row for the order, or ErrNotFound from the store.
func (s *OrderService) GetPaymentForOrder(ctx context.Context, orderID string) (*model.Payment, error) {
	return s.payments.GetPaymentByOrderID(ctx, orderID)
}
