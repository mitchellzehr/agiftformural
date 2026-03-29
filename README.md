# Mural Pay take-home (Go)

SQLite-backed API with `net/http`, deployed to [Railway](https://railway.com/) via **Nixpacks** (no Docker).

HTTP contract: [`openapi.yaml`](openapi.yaml).

## Local run

```bash
export SQLITE_PATH=./app.db
export PORT=8080
# Optional Mural outbound API (payment webhook → CreateTransfer). If unset, a stub is used.
# export MURAL_BASE_URL="https://api-staging.muralpay.com"
# export MURAL_API_KEY="your-staging-key"
go run ./cmd/server
```

The server runs `SeedDefaultProducts` after migrations: three demo catalog rows are inserted: `prod-poster-1`, `prod-stickers-1`, `prod-pin-1`.

## In production

Public base URL (no trailing slash, no `:port` in the browser):

**[https://agiftformural-production.up.railway.app](https://agiftformural-production.up.railway.app)**

Smoke test:

```bash
curl -sS "https://agiftformural-production.up.railway.app/health"
```

Expect `{"ok":true}`.

## Deploy on Railway (GitHub + Nixpacks)

1. Push this repo to GitHub and create a new Railway project from the repository.
2. Railway will detect **Nixpacks**; `nixpacks.toml` builds `./cmd/server` and starts `./server`.

## API examples (curl)

Local server (default `PORT=8080`):

```bash
# Liveness
curl -sS "http://localhost:8080/health"

# Catalog
curl -sS "http://localhost:8080/products"

# Create order (1× poster + 2× sticker sheet → total 54.99 when using seeded catalog)
curl -sS -X POST "http://localhost:8080/orders" \
  -H "Content-Type: application/json" \
  -d '{"currency":"USD","items":[{"product_id":"prod-poster-1","quantity":1},{"product_id":"prod-stickers-1","quantity":2}]}'

# List orders (nested line items)
curl -sS "http://localhost:8080/orders"

# Get one order — substitute the `id` from the create-order response
curl -sS "http://localhost:8080/orders/ORDER_ID"

# Payment webhook — `amount` must match that order's `total_amount` (e.g. 54.99 for the example above)
curl -sS -X POST "http://localhost:8080/webhooks/payment" \
  -H "Content-Type: application/json" \
  -d '{"order_id":"ORDER_ID","amount":54.99}'

# Withdrawals
curl -sS "http://localhost:8080/withdrawals"
curl -sS "http://localhost:8080/withdrawals/WITHDRAWAL_ID"
```

Production (`https://agiftformural-production.up.railway.app`):

```bash
# Liveness
curl -sS "https://agiftformural-production.up.railway.app/health"

# Catalog
curl -sS "https://agiftformural-production.up.railway.app/products"

# Create order (1× poster + 2× sticker sheet → total 54.99 when catalog is seeded)
curl -sS -X POST "https://agiftformural-production.up.railway.app/orders" \
  -H "Content-Type: application/json" \
  -d '{"currency":"USD","items":[{"product_id":"prod-poster-1","quantity":1},{"product_id":"prod-stickers-1","quantity":2}]}'

# List orders
curl -sS "https://agiftformural-production.up.railway.app/orders"

# Get one order — substitute the `id` from the create-order response
curl -sS "https://agiftformural-production.up.railway.app/orders/ORDER_ID"

# Payment webhook — `amount` must match that order's `total_amount`
curl -sS -X POST "https://agiftformural-production.up.railway.app/webhooks/payment" \
  -H "Content-Type: application/json" \
  -d '{"order_id":"ORDER_ID","amount":54.99}'

# Withdrawals
curl -sS "https://agiftformural-production.up.railway.app/withdrawals"
curl -sS "https://agiftformural-production.up.railway.app/withdrawals/WITHDRAWAL_ID"
```

## Current status
- live API with ability

## Future work

- Updrade persistance. Right now it is just SQLite in memory DB. Lost on refresh. This is okay for an exercise but wouldnt hold up in production. We have wrapped this impplementation in interfaces, so swapping SQLite out for something more persistent should be relatively simple
- Mural pay API not called. We have created a client that can be filled in with mural pay information to actually make requests to the sandbox, but I ran out of time here parsing docs


