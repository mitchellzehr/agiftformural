package sqlite

const schema = `
CREATE TABLE IF NOT EXISTS products (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	price REAL NOT NULL,
	created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE TABLE IF NOT EXISTS orders (
	id TEXT PRIMARY KEY,
	status TEXT NOT NULL,
	total_amount REAL NOT NULL,
	currency TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE TABLE IF NOT EXISTS order_items (
	id TEXT PRIMARY KEY,
	order_id TEXT NOT NULL,
	product_id TEXT NOT NULL,
	quantity INTEGER NOT NULL,
	price REAL NOT NULL,
	FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
	FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);

CREATE TABLE IF NOT EXISTS payments (
	id TEXT PRIMARY KEY,
	order_id TEXT NOT NULL,
	amount REAL NOT NULL,
	currency TEXT NOT NULL,
	status TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
	FOREIGN KEY (order_id) REFERENCES orders(id)
);

CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);

CREATE TABLE IF NOT EXISTS withdrawals (
	id TEXT PRIMARY KEY,
	order_id TEXT NOT NULL,
	mural_transfer_id TEXT NOT NULL DEFAULT '',
	amount REAL NOT NULL,
	source_currency TEXT NOT NULL,
	dest_currency TEXT NOT NULL,
	status TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
	FOREIGN KEY (order_id) REFERENCES orders(id)
);

CREATE INDEX IF NOT EXISTS idx_withdrawals_order_id ON withdrawals(order_id);
`
