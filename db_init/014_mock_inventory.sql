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
