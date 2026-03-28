package model

import "time"

type OrderStatus string
type PaymentStatus string
type WithdrawalStatus string
const (
	OrderStatusPendingPayment OrderStatus = "pending_payment"
	OrderStatusPaid           OrderStatus = "paid"
	OrderStatusPayoutPending  OrderStatus = "payout_pending"
	OrderStatusCompleted      OrderStatus = "completed"
	OrderStatusFailed         OrderStatus = "failed"

	PaymentStatusPending  PaymentStatus = "PENDING"
    PaymentStatusReceived PaymentStatus = "RECEIVED"

    WithdrawalStatusPending   WithdrawalStatus = "PENDING"
    WithdrawalStatusCompleted WithdrawalStatus = "COMPLETED"
    WithdrawalStatusFailed    WithdrawalStatus = "FAILED"
)


type Product struct {
    ID        string
    Name      string
    Price     float64
    CreatedAt time.Time
}

type Order struct {
	ID          string
	Status      OrderStatus
	TotalAmount float64
	Currency    string // e.g. "USDC"
	CreatedAt   time.Time
}

type OrderItem struct {
	ID string
	OrderId string
	ProductId string
	Quantity int
	Price float64
}

type Payment struct {
    ID        string
    OrderID   string
    Amount    float64
    Currency  string
    Status    PaymentStatus
    CreatedAt time.Time
}

type Withdrawal struct {
    ID              string
    OrderID         string
    MuralTransferID string
    Amount          float64
    SourceCurrency  string // "USDC"
    DestCurrency    string // "COP"
    Status          WithdrawalStatus
    CreatedAt       time.Time
}