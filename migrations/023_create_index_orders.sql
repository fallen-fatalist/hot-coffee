-- orders idx
CREATE INDEX orders_status_idx ON orders (status);
CREATE INDEX orders_time_idx ON orders (created_at);