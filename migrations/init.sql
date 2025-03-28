
CREATE TABLE customers(
    customer_id SERIAL PRIMARY KEY,
    fullname TEXT NOT NULL,
    phone VARCHAR(20) CONSTRAINT phone_pattern CHECK (phone ~ '^\+?[0-9\-()\s]{7,20}$'),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE orders(
    order_id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    status TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (customer_id) REFERENCES customers (customer_id)
);

CREATE TABLE statuses(
    name TEXT PRIMARY KEY
);


CREATE TABLE order_status_history(
    order_id INTEGER NOT NULL,
    past_status TEXT DEFAULT NULL,
    new_status TEXT NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE
);


-- MUST DO: create index for menu_items_id
CREATE TABLE menu_items(
    menu_item_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT NOT NULL, 
    price NUMERIC NOT NULL CONSTRAINT positive_price CHECK (price > 0)
);

CREATE TABLE order_items(
    menu_item_id INTEGER NOT NULL, 
    order_id INTEGER NOT NULL,
    quantity NUMERIC NOT NULL CONSTRAINT positive_quantity CHECK (quantity > 0),
    customization_info TEXT NOT NULL,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items (menu_item_id) ON DELETE CASCADE,
    FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE price_history(
    menu_item_id INTEGER NOT NULL,
    price_difference INT NOT NULL CONSTRAINT not_the_same CHECK (price_difference != 0),
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (menu_item_id) REFERENCES menu_items (menu_item_id) ON DELETE CASCADE
);



-- MUST DO: Create index for inventory_item_id 
CREATE TABLE inventory(
    inventory_item_id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    price NUMERIC NOT NULL CONSTRAINT positive_price CHECK (price > 0),
    quantity NUMERIC NOT NULL CONSTRAINT positive_quantity CHECK (quantity >= 0),
    unit VARCHAR(20)
);

CREATE TABLE menu_items_ingredients(
    menu_item_id INTEGER NOT NULL,
    inventory_item_id INTEGER NOT NULL,
    quantity NUMERIC NOT NULL CONSTRAINT positive_quantity CHECK (quantity > 0),
    FOREIGN KEY (menu_item_id) REFERENCES menu_items (menu_item_id) ON DELETE CASCADE,
    FOREIGN KEY (inventory_item_id) REFERENCES inventory (inventory_item_id) ON DELETE CASCADE
);

CREATE TABLE units(
    name VARCHAR(20) PRIMARY KEY
);

CREATE TABLE inventory_transactions(
    inventory_item_id INTEGER NOT NULL,
    order_id INTEGER NOT NULL,
    transaction_quantity NUMERIC NOT NULL CONSTRAINT not_the_same CHECK (transaction_quantity != 0),
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (inventory_item_id) REFERENCES inventory (inventory_item_id) ON DELETE CASCADE,
    FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE
);


--Insert mock customers (duplicating existing ones for variety)
INSERT INTO customers (fullname, phone, created_at) VALUES
('John Doe', '+123456789', '2024-02-01 14:23:45 +00:00'),
('Jane Smith', '+987654321', '2023-12-15 08:12:30 +00:00'),
('Alice Johnson', '+112233445', '2023-11-05 17:45:10 +00:00'),
('Charlie Brown', '+443322110', '2023-10-22 21:30:05 +00:00'),
('Eve Davis', '+998877665', '2023-09-14 11:55:20 +00:00'),
('Robert White', '+778899667', '2023-08-27 06:40:33 +00:00'),
('Emma Green', '+667788556', '2023-07-19 15:15:55 +00:00'),
('Oliver Taylor', '+556677445', '2023-06-10 09:05:42 +00:00');


-- Insert mock orders for 2024 (January - December) with 'closed' status
-- Duplicated orders with different months and times
INSERT INTO orders (customer_id, status, created_at) VALUES
(1, 'closed', '2024-01-05 10:20:30'),
(2, 'closed', '2024-02-14 12:30:45'),
(3, 'closed', '2024-03-02 15:40:55'),
(4, 'closed', '2024-04-18 11:10:25'),
(5, 'closed', '2024-05-21 14:15:35'),
(1, 'closed', '2024-06-10 09:05:45'),
(2, 'closed', '2024-07-16 17:30:55'),
(3, 'closed', '2024-08-08 13:20:40'),
(4, 'closed', '2024-09-23 18:25:30'),
(5, 'closed', '2024-10-02 16:45:35'),
(1, 'closed', '2024-11-11 19:15:25'),
(2, 'closed', '2024-12-25 20:30:15'),
(6, 'closed', '2024-01-15 08:30:20'),
(7, 'closed', '2024-02-02 12:10:50'),
(8, 'closed', '2024-03-17 11:25:45'),
(1, 'closed', '2024-04-10 10:05:25'),
(5, 'closed', '2024-05-18 13:10:35'),
(3, 'closed', '2024-06-25 14:30:30'),
(6, 'closed', '2024-07-04 15:40:10'),
(7, 'closed', '2024-08-22 09:25:55'),
(8, 'closed', '2024-09-30 16:10:15'),
(2, 'closed', '2024-10-15 19:40:20'),
(4, 'closed', '2024-11-05 12:50:25'),
(6, 'closed', '2024-12-01 14:35:50');

-- Insert mock orders for 2023 (duplicated and expanded across months)
INSERT INTO orders (customer_id, status, created_at) VALUES
(1, 'closed', '2023-01-15 08:40:20'),
(2, 'closed', '2023-02-20 10:25:35'),
(3, 'closed', '2023-03-05 11:45:10'),
(4, 'closed', '2023-04-10 13:30:55'),
(5, 'closed', '2023-05-17 09:35:50'),
(1, 'closed', '2023-06-08 14:15:25'),
(2, 'closed', '2023-07-09 16:40:30'),
(3, 'closed', '2023-08-22 18:55:15'),
(4, 'closed', '2023-09-14 11:10:25'),
(5, 'closed', '2023-10-05 20:30:35'),
(1, 'closed', '2023-11-13 12:40:10'),
(2, 'closed', '2023-12-19 10:50:55'),
(6, 'closed', '2023-01-20 11:15:30'),
(7, 'closed', '2023-02-04 09:05:10'),
(8, 'closed', '2023-03-11 10:25:25'),
(1, 'closed', '2023-04-22 08:40:35'),
(5, 'closed', '2023-05-25 14:10:50'),
(3, 'closed', '2023-06-12 15:45:25'),
(6, 'closed', '2023-07-20 17:00:30'),
(7, 'closed', '2023-08-30 12:20:35'),
(8, 'closed', '2023-09-10 11:35:40'),
(2, 'closed', '2023-10-22 13:50:10'),
(4, 'closed', '2023-11-01 14:20:50'),
(6, 'closed', '2023-12-10 15:35:15');


-- Adjusted menu item prices
INSERT INTO menu_items (name, description, price) VALUES
('Espresso', 'Strong black coffee made by forcing steam through finely ground coffee beans.', 3.00),
('Latte', 'Espresso with steamed milk and a layer of foam.', 4.50),
('Cappuccino', 'Espresso with steamed milk and foam, topped with chocolate powder.', 4.75),
('Americano', 'Espresso diluted with hot water, creating a smooth and strong coffee.', 3.50),
('Flat White', 'Espresso with steamed milk, very little foam, smoother than a latte.', 4.25),
('Mocha', 'Espresso with chocolate syrup and steamed milk, topped with whipped cream.', 5.50),
('Croissant', 'Flaky, buttery pastry, perfect for a snack.', 3.50),
('Muffin', 'A soft, sweet baked good, with various flavors available.', 3.75),
('Blueberry Muffin', 'Muffin with fresh blueberries inside.', 4.25),
('Chocolate Chip Cookie', 'A sweet and chewy cookie filled with chocolate chips.', 2.50),
('Bagel', 'Boiled and baked dough, served with cream cheese or toppings of your choice.', 3.75),
('Cheesecake', 'Creamy, rich dessert with a graham cracker crust and a smooth topping.', 6.50),
('Tiramisu', 'Italian dessert with layers of coffee-soaked biscuits and mascarpone cream.', 7.00),
('Chocolate Cake', 'Rich and moist cake topped with chocolate frosting.', 5.50),
('Vanilla Cupcake', 'Light and fluffy cake with vanilla frosting on top.', 3.50);


-- TODO: Divide the quantity and information parts in different tables 
-- Insert mock inventory items
INSERT INTO inventory (name, unit, quantity, price) VALUES
('Espresso Beans', 'grams', 10000, 0.04),  -- ~$40 per kg
('Whole Milk', 'liters', 100, 1.20),       -- ~$1.20 per liter
('Almond Milk', 'liters', 50, 2.50),       -- ~$2.50 per liter
('Flour', 'grams', 10000, 0.002),          -- ~$2 per kg
('Sugar', 'grams', 10000, 0.003),          -- ~$3 per kg
('Butter', 'grams', 4000, 0.01),           -- ~$10 per kg
('Blueberries', 'grams', 3000, 0.02),      -- ~$20 per kg
('Chocolate Chips', 'grams', 3000, 0.015), -- ~$15 per kg
('Cream Cheese', 'grams', 2000, 0.012),    -- ~$12 per kg
('Whipped Cream', 'grams', 2000, 0.015),   -- ~$15 per kg
('Vanilla Extract', 'ml', 2000, 0.10),     -- ~$100 per liter
('Graham Crackers', 'grams', 1000, 0.015), -- ~$15 per kg
('Mascarpone Cheese', 'grams', 2000, 0.02),-- ~$20 per kg
('Cocoa Powder', 'grams', 1000, 0.025),    -- ~$25 per kg
('Coffee Syrup', 'ml', 1000, 0.08);        -- ~$80 per liter


-- Insert mock menu item ingredients
INSERT INTO menu_items_ingredients (menu_item_id, inventory_item_id, quantity) VALUES
(1, 1, 18),  -- Espresso: 18g of Espresso Beans per shot
(2, 1, 18),  -- Latte: 18g of Espresso Beans
(2, 2, 0.2),  -- Latte: 200ml of Whole Milk
(3, 1, 18),  -- Cappuccino: 18g of Espresso Beans
(3, 2, 0.150),  -- Cappuccino: 150ml of Whole Milk
(3, 14, 5),  -- Cappuccino: 5g of Cocoa Powder
(4, 1, 18),  -- Americano: 18g of Espresso Beans
(5, 1, 18),  -- Flat White: 18g of Espresso Beans
(5, 2, 0.18),  -- Flat White: 180ml of Whole Milk
(6, 1, 18),  -- Mocha: 18g of Espresso Beans
(6, 2, 0.15),  -- Mocha: 150ml of Whole Milk
(6, 8, 30),  -- Mocha: 30g of Chocolate Chips
(6, 10, 20),  -- Mocha: 20g of Whipped Cream
(7, 4, 100),  -- Croissant: 100g of Flour
(7, 6, 50),  -- Croissant: 50g of Butter
(8, 4, 120),  -- Muffin: 120g of Flour
(8, 5, 30),  -- Muffin: 30g of Sugar
(8, 6, 40),  -- Muffin: 40g of Butter
(9, 4, 120),  -- Blueberry Muffin: 120g of Flour
(9, 5, 30),  -- Blueberry Muffin: 30g of Sugar
(9, 6, 40),  -- Blueberry Muffin: 40g of Butter
(9, 7, 50),  -- Blueberry Muffin: 50g of Blueberries
(10, 8, 25),  -- Chocolate Chip Cookie: 25g of Chocolate Chips
(10, 4, 90),  -- Chocolate Chip Cookie: 90g of Flour
(10, 5, 30),  -- Chocolate Chip Cookie: 30g of Sugar
(10, 6, 20),  -- Chocolate Chip Cookie: 20g of Butter
(11, 4, 150),  -- Bagel: 150g of Flour
(11, 6, 20),  -- Bagel: 20g of Butter
(12, 9, 100),  -- Cheesecake: 100g of Cream Cheese
(12, 5, 50),  -- Cheesecake: 50g of Sugar
(12, 12, 80),  -- Cheesecake: 80g of Graham Crackers
(13, 13, 100),  -- Tiramisu: 100g of Mascarpone Cheese
(13, 5, 30),  -- Tiramisu: 30g of Sugar
(13, 14, 10),  -- Tiramisu: 10g of Cocoa Powder
(13, 15, 50),  -- Tiramisu: 50ml of Coffee Syrup
(14, 8, 50),  -- Chocolate Cake: 50g of Chocolate Chips
(14, 4, 120), -- Chocolate Cake: 120g of Flour
(14, 5, 40),  -- Chocolate Cake: 40g of Sugar
(14, 6, 30),  -- Chocolate Cake: 30g of Butter
(15, 4, 100),  -- Vanilla Cupcake: 100g of Flour
(15, 5, 35),  -- Vanilla Cupcake: 35g of Sugar
(15, 6, 25),  -- Vanilla Cupcake: 25g of Butter
(15, 11, 5);  -- Vanilla Cupcake: 5ml of Vanilla Extract


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

-- Insert mock order items (menu item ID, order ID, quantity, customization info) expanded
INSERT INTO order_items (menu_item_id, order_id, quantity, customization_info) VALUES
(1, 1, 2, 'Extra foam'),
(2, 2, 1, 'No sugar'),
(3, 3, 3, 'With almond milk'),
(1, 4, 1, 'With extra shot'),
(2, 5, 2, 'No blueberries'),
(1, 6, 3, 'Light foam'),
(3, 7, 2, 'With skim milk'),
(1, 8, 1, 'Extra hot'),
(3, 9, 1, 'With extra shot'),
(2, 10, 2, 'No sugar, extra butter'),
(1, 11, 1, 'No foam'),
(3, 12, 1, 'With oat milk'),
(6, 13, 1, 'Hotter than usual'),
(7, 14, 2, 'Extra cream'),
(8, 15, 1, 'No nuts'),
(1, 16, 2, 'Less foam'),
(3, 17, 2, 'With extra vanilla'),
(2, 18, 1, 'More blueberries'),
(7, 19, 1, 'Extra hot'),
(8, 20, 3, 'With extra sugar'),
(5, 21, 1, 'No frosting'),
(3, 22, 2, 'With coconut milk'),
(2, 23, 1, 'No sugar, extra honey'),
(6, 24, 2, 'With extra soy milk');

-- Insert mock statuses (open, In Progress, closed, Cancelled)
-- Keep the same status logic (unchanged)
INSERT INTO statuses (name) VALUES
('open'),
('in progress'),
('rejected'),
('closed');

-- Insert mock order status history (expanded with more orders)
INSERT INTO order_status_history (order_id, past_status, new_status) VALUES
(2, 'open', 'closed'),
(3, 'open', 'closed'),
(4, 'in progress', 'closed'),
(5, 'open', 'closed'),
(6, 'in progress', 'closed'),
(9, 'in progress', 'closed'),
(10, 'open', 'closed'),
(11, 'in progress', 'closed'),
(12, 'closed', 'closed'),
(13, 'open', 'closed'),
(14, 'closed', 'in progress'),
(16, 'in progress', 'closed'),
(18, 'open', 'closed'),
(19, 'in progress', 'closed'),
(20, 'closed', 'closed'),
(22, 'closed', 'in progress'),
(23, 'open', 'closed'),
(24, 'closed', 'closed');


-- Insert mock price history (expanded with more items)
INSERT INTO price_history (menu_item_id, price_difference) VALUES
(2, 1),
(3, 2),
(4, 1),  -- New menu item
(5, 1);  -- New menu item

-- Insert mock inventory transactions (expanded with more items used)
INSERT INTO inventory_transactions (inventory_item_id, order_id, transaction_quantity) VALUES
(1, 1, -40),  -- 40 Espresso Shots used for orders
(2, 1, -2000), -- 2000ml of milk used for orders
(3, 1, -1000),  -- 1000g of flour used for orders
(4, 1, -300),  -- 300g of blueberries used for orders
(5, 1, 600);   -- 600g of sugar restocked

-- orders idx
CREATE INDEX orders_status_idx ON orders (status);
CREATE INDEX orders_time_idx ON orders (created_at);
--order_items as oi idx
CREATE INDEX oi_quantity_idx ON order_items (quantity);
-- indexes for foreign key in junction table
CREATE INDEX oi_menuid_idx ON order_items (menu_item_id);
CREATE INDEX oi_orderid_idx ON order_items (order_id);
-- menu_items_ingredients as mii
CREATE INDEX mii_menu_item_id ON menu_items_ingredients (menu_item_id);
CREATE INDEX mii_inventory_item_id ON menu_items_ingredients (inventory_item_id);
    