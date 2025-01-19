-- CREATE TYPE status_type AS ENUM ('accepted', 'rejected');
-- CREATE TABLE unit_type AS ENUM ('shots', 'ml', 'g');


CREATE TABLE customers(
    customer_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    phone VARCHAR(20) CHECK (phone ~ '^\+?[0-9\-()\s]{7,20}$')
);

CREATE TABLE orders(
    order_id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (customer_id) REFERENCES customers (customer_id)
);

CREATE TABLE statuses(
    status_id SERIAL PRIMARY KEY,
    name TEXT
);


CREATE TABLE order_status_history(
    order_id INTEGER NOT NULL,
    past_status TEXT NOT NULL,
    new_status TEXT NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (order_id) REFERENCES orders (order_id) 
);


CREATE TABLE menu_items(
    menu_item_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT NOT NULL, 
    price INT NOT NULL
);

CREATE TABLE order_items(
    menu_item_id INTEGER NOT NULL, 
    order_id INTEGER NOT NULL,
    quantity DECIMAL(10, 5) NOT NULL,
    customization_info TEXT NOT NULL,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items (menu_item_id),
    FOREIGN KEY (order_id) REFERENCES orders (order_id)
);

CREATE TABLE price_history(
    menu_item_id INTEGER NOT NULL,
    price_difference INT NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (menu_item_id) REFERENCES menu_items (menu_item_id)
);



CREATE TABLE inventory(
    inventory_item_id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    quantity DECIMAL(10,5) NOT NULL,
    unit VARCHAR(20)
);

CREATE TABLE menu_items_ingredients(
    menu_item_id INTEGER NOT NULL,
    inventory_item_id INTEGER NOT NULL,
    quantity DECIMAL(10, 5) NOT NULL,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items (menu_item_id),
    FOREIGN KEY (inventory_item_id) REFERENCES inventory (inventory_item_id)
);

CREATE TABLE units(
    unit_id SERIAL PRIMARY KEY,
    name VARCHAR(20)
);

CREATE TABLE inventory_transactions(
    inventory_item_id SERIAL PRIMARY KEY,
    transaction_quantity DECIMAL(10, 5) NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (inventory_item_id) REFERENCES inventory (inventory_item_id)
);