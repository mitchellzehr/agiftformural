# Mural Pay take-home (Go)

SQLite-backed API with `net/http`, deployed to [Railway](https://railway.com/) via **Nixpacks** (no Docker).

HTTP contract: [`openapi.yaml`](openapi.yaml).

## Local run

```bash
export MURAL_API_KEY="your-staging-api-key"
export SQLITE_PATH=./app.db
export PORT=8080
go run ./cmd/server
```

## In production

Public base URL (no trailing slash, no `:port` in the browser):

**[https://agiftformural-production.up.railway.app](https://agiftformural-production.up.railway.app)**

Smoke test:

```bash
curl -sS "https://agiftformural-production.up.railway.app/health"
```

Expect `{"ok":true}`. Use the same host as `BASE` in the curl examples below.

## Deploy on Railway (GitHub + Nixpacks)

1. Push this repo to GitHub and create a new Railway project from the repository.
2. Railway will detect **Nixpacks**; `nixpacks.toml` builds `./cmd/server` and starts `./server`.

## API spec (quick reference)

JSON only for bodies (`Content-Type: application/json`). Request bodies must not include unknown fields. Timestamps are RFC 3339 strings (server may use nanosecond precision). Canonical detail lives in [`openapi.yaml`](openapi.yaml).

### Endpoints

| Method | Path | Summary |
| --- | --- | --- |
| `GET` | `/health` | Liveness |
| `GET` | `/products` | List catalog products |
| `GET` | `/orders` | List orders, each with nested line items |
| `POST` | `/orders` | Create order in `pending_payment`; line unit prices from catalog |
| `GET` | `/orders/{id}` | Order + line items |
| `POST` | `/webhooks/payment` | Record payment when `amount` matches order total (float tolerance ~1e-6); body max 1 MiB |
| `GET` | `/withdrawals` | List withdrawals |
| `GET` | `/withdrawals/{id}` | Single withdrawal |

### Request bodies

**`POST /orders`** (`CreateOrderRequest`)

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `currency` | string | yes | e.g. `USD` / `USDC` |
| `id` | string | no | Order id; server generates UUID if omitted |
| `items` | array | no | Omit or `[]` for an order with no lines |
| `items[].product_id` | string | yes (per line) | Must exist in catalog |
| `items[].quantity` | integer | yes (per line) | ≥ 1 |
| `items[].id` | string | no | Line id; server UUID if omitted |

**`POST /webhooks/payment`** (`PaymentWebhookRequest`)

| Field | Type | Required |
| --- | --- | --- |
| `order_id` | string | yes |
| `amount` | number | yes | Must equal the order’s `total_amount` within tolerance |

### Successful response shapes

| Call | Status | Body |
| --- | --- | --- |
| `GET /health` | 200 | `{ "ok": true }` |
| `GET /products` | 200 | `{ "products": [ { "id", "name", "price", "created_at" } ] }` |
| `GET /orders` | 200 | `{ "orders": [ { "order": Order, "items": OrderItem[] } ] }` |
| `POST /orders` | 201 | `{ "order": Order, "items": OrderItem[] }` |
| `GET /orders/{id}` | 200 | `{ "order": Order, "items": OrderItem[] }` |
| `POST /webhooks/payment` | 200 | `{ "ok": true }` |
| `GET /withdrawals` | 200 | `{ "withdrawals": Withdrawal[] }` |
| `GET /withdrawals/{id}` | 200 | `{ "withdrawal": Withdrawal }` |

### Shared types

- **Order:** `id`, `status`, `total_amount`, `currency`, `created_at`
- **Order `status`:** `pending_payment` · `paid` · `payout_pending` · `completed` · `failed`
- **OrderItem:** `id`, `order_id`, `product_id`, `quantity`, `price` (unit price at order time)
- **Withdrawal:** `id`, `order_id`, `mural_transfer_id`, `amount`, `source_currency`, `dest_currency`, `status`, `created_at`
- **Withdrawal `status`:** `PENDING` · `COMPLETED` · `FAILED`
- **Errors (typical):** `{ "error": "…" }` on 4xx / 5xx

## API examples (curl)

Set `BASE` to where the server listens:

- **Local:** `http://localhost:8080` (or your `PORT`)
- **Production:** `https://agiftformural-production.up.railway.app`

```bash
BASE="http://localhost:8080"

# Liveness
curl -sS "$BASE/health"

# Catalog (use returned product ids for creating orders)
curl -sS "$BASE/products"

# Create order — replace PRODUCT_ID with a value from GET /products
curl -sS -X POST "$BASE/orders" \
  -H "Content-Type: application/json" \
  -d '{"currency":"USD","items":[{"product_id":"PRODUCT_ID","quantity":1}]}'

# List orders (nested line items)
curl -sS "$BASE/orders"

# Get one order — replace ORDER_ID
curl -sS "$BASE/orders/ORDER_ID"

# Payment webhook — `amount` must match the order's `total_amount` from GET /orders/{id}
curl -sS -X POST "$BASE/webhooks/payment" \
  -H "Content-Type: application/json" \
  -d '{"order_id":"ORDER_ID","amount":99.5}'

# Withdrawals
curl -sS "$BASE/withdrawals"
curl -sS "$BASE/withdrawals/WITHDRAWAL_ID"
```

## Current status

- Live HTTP API (orders, products, withdrawals, payment webhook). See **API spec** above and [`openapi.yaml`](openapi.yaml).

## Future work

- Updrade persistance. Right now it is just SQLite in memory DB. Lost on refresh. This is okay for an exercise but wouldnt hold up in production. We have wrapped this impplementation in interfaces, so swapping SQLite out for something more persistent should be relatively simple
- Mural pay API not called. We have created a client that can be filled in with mural pay information to actually make requests to the sandbox, but I ran out of time here parsing docs


