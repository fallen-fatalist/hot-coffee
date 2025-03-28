CREATE TABLE order_status_history(
    order_id INTEGER NOT NULL,
    past_status TEXT DEFAULT NULL,
    new_status TEXT NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE
);