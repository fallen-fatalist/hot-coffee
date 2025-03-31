CREATE TABLE order_items(
    menu_item_id INTEGER NOT NULL, 
    order_id INTEGER NOT NULL,
    quantity NUMERIC NOT NULL CONSTRAINT positive_quantity CHECK (quantity > 0),
    customization_info TEXT NOT NULL,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items (menu_item_id) ON DELETE CASCADE,
    FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);