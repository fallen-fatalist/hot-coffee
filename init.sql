-- CREATE TYPE status_type AS ENUM ('accepted', 'rejected');
-- CREATE TABLE unit_type AS ENUM ('shots', 'ml', 'g');

-----------------
-- Orders part --
-----------------

CREATE TABLE orders(
    order_id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL ,
    status_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers (customer_id),
    FOREIGN KEY (status_id) REFERENCES statuses (status_id)
);

CREATE TABLE statuses(
    status_id SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE customers(
    customer_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    phone VARCHAR(20) CHECK (phone ~ '^\+?[0-9\-()\s]{7,20}$')
);

CREATE TABLE order_status_history(
    order_id INTEGER NOT NULL,
    past_status INTEGER NOT NULL,
    new_status INTEGER NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders (order_id),
    FOREIGN KEY (past_status) REFERENCES statuses (status_id),
    FOREIGN KEY (new_status) REFERENCES STATUSES (status_id)
);

CREATE TABLE order_items(
    menu_item_id INTEGER NOT NULL, 
    order_id INTEGER NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    customization_info TEXT NOT NULL,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items (menu_item_id),
    FOREIGN KEY (order_id) REFERENCES orders (order_id)
);

---------------
-- Menu part --
---------------

CREATE TABLE menu_items(
    menu_item_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT NOT NULL, 
    price INT NOT NULL,
    created_at TIMESTAMPTZ
);

CREATE TABLE price_history(
    menu_item_id INTEGER NOT NULL,
    price_difference INT NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items (menu_item_id)
);

CREATE TABLE menu_items_ingredients(
    menu_item_id INTEGER NOT NULL,
    inventory_item_id INTEGER NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items (menu_item_id),
    FOREIGN KEY (inventory_item_id) REFERENCES inventory (inventory_item_id)
);

--------------------
-- Inventory part --
--------------------

CREATE TABLE inventory(
    inventory_item_id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    unit_id INT NOT NULL,
    created_at TIMESTAMPTZ,
    FOREIGN KEY (unit_id) REFERENCES units (unit_id)
);

CREATE TABLE units(
    unit_id SERIAL PRIMARY KEY,
    name VARCHAR(20)
);

CREATE TABLE inventory_transactions(
    inventory_item_id SERIAL PRIMARY KEY,
    transaction_quantity DECIMAL(10, 2) NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (inventory_item_id) REFERENCES inventory (inventory_item_id)
);