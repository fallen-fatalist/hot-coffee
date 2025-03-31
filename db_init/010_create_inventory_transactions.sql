CREATE TABLE inventory_transactions(
    inventory_item_id INTEGER NOT NULL,
    order_id INTEGER NOT NULL,
    transaction_quantity NUMERIC NOT NULL CONSTRAINT not_the_same CHECK (transaction_quantity != 0),
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (inventory_item_id) REFERENCES inventory (inventory_item_id) ON DELETE CASCADE,
    FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE
);