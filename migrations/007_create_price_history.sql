CREATE TABLE price_history(
    menu_item_id INTEGER NOT NULL,
    price_difference INT NOT NULL CONSTRAINT not_the_same CHECK (price_difference != 0),
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (menu_item_id) REFERENCES menu_items (menu_item_id) ON DELETE CASCADE
);
