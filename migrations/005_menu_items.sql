
-- MUST DO: create index for menu_items_id
CREATE TABLE menu_items(
    menu_item_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT NOT NULL, 
    price NUMERIC NOT NULL CONSTRAINT positive_price CHECK (price > 0)
);