CREATE TYPE status_type AS ENUM ('accepted', 'rejected');

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
    FOREIGN KEY (order_id) REFERENCES orders (order_id)   

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

