package handler

import (
	"context"
	"net/http"
	"time"

	"mural/internal/model"
)

// ProductCatalog is what ListProducts needs from the store.
type ProductCatalog interface {
	ListProducts(ctx context.Context) ([]model.Product, error)
}

type productResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
}

func (s *Server) ListProducts(w http.ResponseWriter, r *http.Request) {
	items, err := s.products.ListProducts(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	out := make([]productResponse, 0, len(items))
	for _, p := range items {
		out = append(out, productResponse{
			ID:        p.ID,
			Name:      p.Name,
			Price:     p.Price,
			CreatedAt: p.CreatedAt.UTC().Format(time.RFC3339Nano),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"products": out})
}
