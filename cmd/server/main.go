package main

import (
	"context"
	"log"
	"net/http"

	"mural/internal/config"
	"mural/internal/handler"
	sqliteRepo "mural/internal/repo/sqlite"
	"mural/internal/service"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db, err := sqliteRepo.Open(cfg.SQLitePath)
	if err != nil {
		log.Fatalf("sqlite open %q: %v", cfg.SQLitePath, err)
	}
	defer db.Close()

	if err := sqliteRepo.InitSchema(ctx, db); err != nil {
		log.Fatalf("schema: %v", err)
	}

	repos := sqliteRepo.NewRepos(db)
	store := &repos
	orderSvc := service.NewOrderService(store, store, store)
	h := handler.New(handler.Deps{
		Products:    store,
		OrderSvc:    orderSvc,
		Withdrawals: store,
	})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /products", h.ListProducts)
	mux.HandleFunc("POST /orders", h.CreateOrder)
	mux.HandleFunc("GET /orders", h.ListOrders)
	mux.HandleFunc("GET /orders/{id}", h.GetOrder)
	mux.HandleFunc("POST /webhooks/payment", h.PaymentWebhook)
	mux.HandleFunc("GET /withdrawals", h.ListWithdrawals)
	mux.HandleFunc("GET /withdrawals/{id}", h.GetWithdrawal)

	addr := ":" + cfg.Port
	log.Printf("listening on %s (SQLITE_PATH=%s)", addr, cfg.SQLitePath)
	log.Fatal(http.ListenAndServe(addr, mux))
}
