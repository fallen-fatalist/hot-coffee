-- enum type for status
CREATE TYPE status AS ENUM ('open', 'in progress', 'rejected', 'closed');

CREATE TABLE orders(
    order_id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    status status,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (customer_id) REFERENCES customers (customer_id)
);