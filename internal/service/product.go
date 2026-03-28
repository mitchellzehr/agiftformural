package service

import (
	"context"

	"mural/internal/model"
)

// ProductService exposes catalog operations backed by persistence.
type ProductService struct {
	products ProductReader
}

// NewProductService builds a catalog service over a ProductReader (e.g. repos.Products).
func NewProductService(products ProductReader) *ProductService {
	return &ProductService{products: products}
}

// ListProducts returns the full catalog, ordered by the persistence layer.
func (s *ProductService) ListProducts(ctx context.Context) ([]model.Product, error) {
	return s.products.ListProducts(ctx)
}

// GetProduct returns a single product or an error (including ErrNotFound from the store).
func (s *ProductService) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	return s.products.GetProduct(ctx, id)
}
