package service

import (
	"context"

	"mural/internal/model"
)

// ProductReader supports GET /products (and checkout validation).
type ProductReader interface {
	ListProducts(ctx context.Context) ([]model.Product, error)
	GetProduct(ctx context.Context, id string) (*model.Product, error)
}

// OrderStore supports POST /orders, GET /orders, GET /orders/{id}.
type OrderStore interface {
	CreateOrderWithItems(ctx context.Context, o *model.Order, items []model.OrderItem) error
	GetOrderByID(ctx context.Context, id string) (*model.Order, []model.OrderItem, error)
	ListOrders(ctx context.Context) ([]model.Order, []model.OrderItem, error)
}

// PaymentStore supports inserting payment rows (append-only).
type PaymentStore interface {
	CreatePayment(ctx context.Context, p *model.Payment) error
	GetPaymentByOrderID(ctx context.Context, orderID string) (*model.Payment, error)
}

// WithdrawalStore supports listing withdrawals and append-only creates.
type WithdrawalStore interface {
	ListWithdrawals(ctx context.Context) ([]model.Withdrawal, error)
	GetWithdrawal(ctx context.Context, id string) (*model.Withdrawal, error)
	CreateWithdrawal(ctx context.Context, w *model.Withdrawal) error
}

// Repos groups repository interfaces for wiring in main and persistence constructors.
type Repos struct {
	Products    ProductReader
	Orders      OrderStore
	Payments    PaymentStore
	Withdrawals WithdrawalStore
}
