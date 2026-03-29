package main

import (
	"context"
	"log"
	"net/http"

	"mural/internal/config"
	"mural/internal/handler"
	"mural/internal/mural"
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
	if err := sqliteRepo.SeedDefaultProducts(ctx, db); err != nil {
		log.Fatalf("seed products: %v", err)
	}

	repos := sqliteRepo.NewRepos(db)
	store := &repos

	var muralSvc service.MuralClient = mural.NewClient(cfg.MuralURL, cfg.MuralAPIKey)
	if cfg.MuralAPIKey != "" {
		muralSvc = mural.NewClient(cfg.MuralURL, cfg.MuralAPIKey)
	} else {
		log.Print("mural: MURAL_API_KEY unset; using StubClient (no outbound CreateTransfer HTTP)")
	}

	orderSvc := service.NewOrderService(store, store, store, store, muralSvc)
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
