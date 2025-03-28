CREATE TABLE customers(
    customer_id SERIAL PRIMARY KEY,
    fullname TEXT NOT NULL,
    phone VARCHAR(20) CONSTRAINT phone_pattern CHECK (phone ~ '^\+?[0-9\-()\s]{7,20}$'),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
