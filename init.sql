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
    price DECIMAL(10,2) NOT NULL
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
    price DECIMAL(10, 2) NOT NULL,
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

-- TEMPORARY MOCK DATA INSERTS

-- Insert mock inventory data
INSERT INTO inventory (name, price, quantity, unit) VALUES
('Espresso Shot', 300, 500, 'shots'),
('Milk', 200, 5000, 'ml'),
('Flour', 100, 10000, 'g'),
('Blueberries', 50, 2000, 'g'),
('Sugar', 10, 5000, 'g');

-- Insert mock customers
INSERT INTO customers (name, surname, phone) VALUES
('John', 'Doe', '+123456789'),
('Jane', 'Smith', '+987654321'),
('Alice', 'Johnson', '+112233445');

-- Insert mock menu items
INSERT INTO menu_items (name, description, price) VALUES
('Cappuccino', 'Espresso with steamed milk foam', 5),
('Blueberry Muffin', 'Freshly baked with blueberries', 3),
('Latte', 'Espresso with steamed milk', 4);

-- Insert mock menu_items_ingredients linking menu items to inventory
INSERT INTO menu_items_ingredients (menu_item_id, inventory_item_id, quantity) VALUES
(1, 1, 1.0), -- Cappuccino uses 1 shot of espresso
(1, 2, 200.0), -- Cappuccino uses 200ml of milk
(2, 3, 100.0), -- Blueberry Muffin uses 100g of flour
(2, 4, 50.0), -- Blueberry Muffin uses 50g of blueberries
(3, 1, 1.0), -- Latte uses 1 shot of espresso
(3, 2, 300.0); -- Latte uses 300ml of milk

-- Insert mock orders
INSERT INTO orders (customer_id, status) VALUES
(1, 'accepted'),
(2, 'in_progress'),
(3, 'completed');

-- Insert mock order items
INSERT INTO order_items (menu_item_id, order_id, quantity, customization_info) VALUES
(1, 1, 1, 'Extra foam'),
(2, 2, 2, 'No sugar'),
(3, 3, 1, 'With almond milk');

-- Insert mock statuses
INSERT INTO statuses (name) VALUES
('accepted'),
('in_progress'),
('completed'),
('cancelled');

-- Insert mock order status history
INSERT INTO order_status_history (order_id, past_status, new_status) VALUES
(1, 'pending', 'accepted'),
(2, 'accepted', 'in_progress'),
(3, 'in_progress', 'completed');

-- Insert mock price history
INSERT INTO price_history (menu_item_id, price_difference) VALUES
(1, 1), -- Cappuccino price increased by $1
(2, 0), -- Blueberry Muffin price unchanged
(3, 2); -- Latte price increased by $2

-- Insert mock inventory transactions
INSERT INTO inventory_transactions (inventory_item_id, transaction_quantity) VALUES
(1, -10), -- 10 Espresso Shots used
(2, -500), -- 500ml of milk used
(3, -300), -- 300g of flour used
(4, -100), -- 100g of blueberries used
(5, 500); -- 500g of sugar restocked
