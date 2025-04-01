-- MUST DO: Create index for inventory_item_id 
CREATE TYPE unit AS ENUM ('grams', 'liters', 'ml', 'kg');

CREATE TABLE inventory(
    inventory_item_id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    price NUMERIC NOT NULL CONSTRAINT positive_price CHECK (price > 0),
    quantity NUMERIC NOT NULL CONSTRAINT positive_quantity CHECK (quantity >= 0),
    unit unit
);