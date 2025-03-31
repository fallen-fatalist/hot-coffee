CREATE TABLE menu_items_ingredients(
    menu_item_id INTEGER NOT NULL,
    inventory_item_id INTEGER NOT NULL,
    quantity NUMERIC NOT NULL CONSTRAINT positive_quantity CHECK (quantity > 0),
    FOREIGN KEY (menu_item_id) REFERENCES menu_items (menu_item_id) ON DELETE CASCADE,
    FOREIGN KEY (inventory_item_id) REFERENCES inventory (inventory_item_id) ON DELETE CASCADE
);