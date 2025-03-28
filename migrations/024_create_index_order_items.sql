--order_items as oi idx
CREATE INDEX oi_quantity_idx ON order_items (quantity);
CREATE INDEX oi_orderid_idx ON order_items (order_id);