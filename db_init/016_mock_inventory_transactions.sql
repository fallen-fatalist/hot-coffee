-- Insert mock inventory transactions (usage and restock events)
INSERT INTO inventory_transactions (inventory_item_id, order_id, transaction_quantity) VALUES
(1, 1, -40),  -- 40 grams of Espresso Beans used
(2, 1, -2000),  -- 2000 ml of Whole Milk used
(3, 1, -800),  -- 800 ml of Almond Milk used
(4, 1, -1000),  -- 1000 grams of Flour used
(5, 1, -600),   -- 600 grams of Sugar used
(6, 1, -300),   -- 300 grams of Butter used
(7, 1, -300),   -- 300 grams of Blueberries used
(8, 1, -1000),  -- 1000 grams of Chocolate Chips used
(9, 1, -500),   -- 500 grams of Cream Cheese used
(10, 1, -200),  -- 200 grams of Whipped Cream used
(11, 1, -100),  -- 100 ml of Vanilla Extract used
(12, 1, -500),  -- 500 grams of Graham Crackers used
(13, 1, -200),  -- 200 grams of Mascarpone Cheese used
(14, 1, -100),  -- 100 grams of Cocoa Powder used
(15, 1, -200);  -- 200 ml of Coffee Syrup used

-- Restocking inventory items
INSERT INTO inventory_transactions (inventory_item_id, order_id, transaction_quantity) VALUES
(1, 2, 1000),  -- Restocking Espresso Beans
(2, 2, 2000),  -- Restocking Whole Milk
(3, 2, 1000),  -- Restocking Almond Milk
(4, 2, 3000),  -- Restocking Flour
(5, 2, 1000),  -- Restocking Sugar
(6, 2, 500),   -- Restocking Butter
(7, 2, 500),   -- Restocking Blueberries
(8, 2, 1500),  -- Restocking Chocolate Chips
(9, 2, 700),   -- Restocking Cream Cheese
(10, 2, 300),  -- Restocking Whipped Cream
(11, 2, 300),  -- Restocking Vanilla Extract
(12, 2, 1500), -- Restocking Graham Crackers
(13, 2, 700),  -- Restocking Mascarpone Cheese
(14, 2, 400),  -- Restocking Cocoa Powder
(15, 2, 300);  -- Restocking Coffee Syrup