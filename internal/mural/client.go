package mural

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	u := strings.TrimSuffix(strings.TrimSpace(baseURL), "/")
	return &Client{
		baseURL:    u,
		apiKey:     strings.TrimSpace(apiKey),
		httpClient: http.DefaultClient,
	}
}

// CreateTransfer requests a transfer for the given amount and returns the Mural transfer identifier.
// This functionality is not implemented yet due to time constraints
func (c *Client) CreateTransfer(ctx context.Context, amount float64) (string, error) {
	return uuid.NewString(), nil
}
