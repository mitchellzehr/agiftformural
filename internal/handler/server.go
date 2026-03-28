package handler

import "mural/internal/service"

// Server is the HTTP surface; dependencies are narrow interfaces declared alongside handlers.
type Server struct {
	products    ProductCatalog
	orderSvc    service.OrderService
	withdrawals WithdrawalReader
}

// Deps bundles constructor inputs for New.
type Deps struct {
	Products    ProductCatalog
	OrderSvc    service.OrderService
	Withdrawals WithdrawalReader
}

// New returns a Server wired with the given dependencies.
func New(d Deps) *Server {
	return &Server{
		products:    d.Products,
		orderSvc:    d.OrderSvc,
		withdrawals: d.Withdrawals,
	}
}
