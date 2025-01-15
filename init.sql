CREATE TYPE status_type AS ENUM ('accepted', 'rejected');
CREATE TABLE unit_type AS ENUM ('shots', 'ml', 'g');

CREATE TABLE orders(
    order_id INT PRIMARY KEY,
    customer_id INT NOT NULL,
    status status_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers (customer_id)
);

CREATE TABLE order_items(
    menu_item_id INT NOT NULL, 
    order_id INT NOT NULL,
    quantity INT NOT NULL,
    customiization_info TEXT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders (order_id),
    FOREIGN KEY (menu_item_id) REFERENCES menu_items (menu_item_id) 
);

CREATE TABLE order_status_history(
    order_id INT NOT NULL,
    past_status status_type NOT NULL,
    new_status status_type NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders (order_id)
);

CREATE TABLE customers(
    customer_id INT PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    contact_information TEXT NOT NULL
);

CREATE TABLE menu_items(
    menu_item_id INT PRIMARY KEY,
    menu_items_ingredients_id INT NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT NOT NULL, 
    price INT NOT NULL,
    FOREIGN KEY (menu_items_ingredients_id) REFERENCES menu_items_ingredients (menu_items_ingredients_id) 
);

CREATE TABLE menu_items_ingredients(
    menu_items_ingredients_id INT PRIMARY KEY,
    inventory_item_id INT NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (inventory_item_id) REFERENCES inventory (inventory_item_id)
);

CREATE TABLE price_history(
    menu_item_id INT NOT NULL,
    price_difference INT NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items (menu_item_id)
);

CREATE TABLE inventory(
    inventory_item_id INT PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    quantity INT NOT NULL,
    unit unit_type NOT NULL
);

CREATE TABLE inventory_transactions(
    inventory_item_id INT NOT NULL,
    transaction_quantity INT NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (inventory_item_id) REFERENCES inventory (inventory_item_id)
);