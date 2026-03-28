package mural

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client calls the Mural HTTP API. CreateTransfer matches service.MuralClient for wiring into OrderService.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient returns a client for baseURL (e.g. https://api-staging.muralpay.com). Trailing slashes are trimmed.
// apiKey is sent as Authorization: Bearer when non-empty.
func NewClient(baseURL, apiKey string) *Client {
	u := strings.TrimSuffix(strings.TrimSpace(baseURL), "/")
	return &Client{
		baseURL:    u,
		apiKey:     strings.TrimSpace(apiKey),
		httpClient: http.DefaultClient,
	}
}

// CreateTransfer requests a transfer for the given amount and returns the Mural transfer identifier.
// The request path and JSON shape are aligned with typical REST patterns; adjust when wiring to the live Mural spec.
func (c *Client) CreateTransfer(ctx context.Context, amount float64) (string, error) {
	if c.baseURL == "" {
		return "", fmt.Errorf("mural: empty base URL")
	}
	body, err := json.Marshal(map[string]any{
		"amount": amount,
	})
	if err != nil {
		return "", err
	}
	url := c.baseURL + "/transfers"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("mural: POST %s: %s: %s", url, resp.Status, bytes.TrimSpace(respBody))
	}
	var out struct {
		ID         string `json:"id"`
		TransferID string `json:"transfer_id"`
	}
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", fmt.Errorf("mural: decode response: %w", err)
	}
	if out.ID != "" {
		return out.ID, nil
	}
	if out.TransferID != "" {
		return out.TransferID, nil
	}
	return "", fmt.Errorf("mural: response missing transfer id")
}
