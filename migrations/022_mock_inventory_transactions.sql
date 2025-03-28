-- Insert mock inventory transactions (expanded with more items used)
INSERT INTO inventory_transactions (inventory_item_id, order_id, transaction_quantity) VALUES
(1, 1, -40),  -- 40 Espresso Shots used for orders
(2, 1, -2000), -- 2000ml of milk used for orders
(3, 1, -1000),  -- 1000g of flour used for orders
(4, 1, -300),  -- 300g of blueberries used for orders
(5, 1, 600);   -- 600g of sugar restocked