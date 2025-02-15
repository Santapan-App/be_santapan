CREATE TABLE food_nutrition (
    id BIGSERIAL PRIMARY KEY,
    food_name VARCHAR(100) NOT NULL,
    calories INT NOT NULL,
    protein INT NOT NULL,
    fat INT NOT NULL,
    carbohydrates INT NOT NULL,
    sugar INT NOT NULL
);

CREATE TABLE menu_nutrition (
    id BIGSERIAL PRIMARY KEY,
    menu_id BIGINT REFERENCES menu(id) ON DELETE CASCADE,
    food_id BIGINT REFERENCES food_nutrition(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on food_name for faster searching
CREATE INDEX idx_food_name ON food_nutrition(food_name);

INSERT INTO food_nutrition (food_name, calories, protein, fat, carbohydrates, sugar) VALUES
('apple_pie', 320, 2, 14, 48, 24),
('baby_back_ribs', 300, 20, 22, 8, 1),
('baklava', 330, 6, 18, 38, 25),
('beef_carpaccio', 150, 20, 7, 1, 0),
('beef_tartare', 210, 23, 12, 0, 0),
('beet_salad', 150, 5, 9, 14, 4),
('beignets', 350, 5, 20, 38, 15),
('bibimbap', 500, 20, 15, 70, 5),
('bread_pudding', 400, 8, 12, 64, 22),
( 'breakfast_burrito', 450, 25, 20, 45, 6),
( 'bruschetta', 200, 5, 8, 28, 2),
( 'caesar_salad', 300, 10, 20, 12, 2),
( 'cannoli', 250, 5, 12, 34, 20),
( 'caprese_salad', 220, 10, 15, 9, 3),
( 'carrot_cake', 400, 5, 20, 60, 35),
( 'ceviche', 180, 22, 6, 4, 0),
( 'cheese_plate', 450, 25, 35, 15, 1),
( 'cheesecake', 400, 7, 28, 30, 25),
( 'chicken_curry', 500, 25, 20, 30, 5),
( 'chicken_quesadilla', 450, 30, 20, 35, 3),
( 'chicken_wings', 400, 30, 30, 0, 0),
( 'chocolate_cake', 350, 5, 18, 52, 32),
( 'chocolate_mousse', 250, 4, 18, 30, 20),
( 'churros', 400, 6, 20, 54, 24),
( 'clam_chowder', 350, 12, 20, 30, 1),
( 'club_sandwich', 600, 30, 28, 40, 4),
( 'crab_cakes', 250, 15, 14, 18, 0),
( 'creme_brulee', 300, 6, 20, 28, 28),
( 'croque_madame', 450, 25, 25, 40, 5),
( 'cup_cakes', 250, 3, 10, 40, 18),
( 'deviled_eggs', 150, 12, 10, 2, 0),
( 'donuts', 300, 4, 18, 35, 15),
( 'dumplings', 200, 6, 4, 36, 1),
( 'edamame', 120, 11, 5, 9, 0),
( 'eggs_benedict', 500, 20, 30, 30, 1),
( 'escargots', 250, 20, 18, 0, 0),
( 'falafel', 350, 15, 20, 45, 2),
( 'filet_mignon', 600, 45, 40, 0, 0),
( 'fish_and_chips', 700, 25, 30, 80, 1),
( 'foie_gras', 500, 10, 45, 0, 0),
( 'french_fries', 365, 3, 17, 63, 0),
( 'french_onion_soup', 250, 8, 15, 30, 1),
( 'french_toast', 350, 10, 16, 45, 15),
( 'fried_calamari', 400, 10, 20, 35, 0),
( 'fried_rice', 350, 8, 10, 60, 1),
( 'frozen_yogurt', 150, 5, 2, 30, 20),
( 'garlic_bread', 300, 8, 14, 36, 1),
( 'gnocchi', 250, 7, 1, 53, 1),
( 'greek_salad', 200, 5, 12, 9, 2),
( 'grilled_cheese_sandwich', 400, 12, 24, 36, 4),
( 'grilled_salmon', 350, 40, 20, 0, 0),
( 'guacamole', 150, 2, 15, 8, 1),
( 'gyoza', 220, 7, 8, 30, 1),
( 'hamburger', 500, 30, 25, 40, 7),
( 'hot_and_sour_soup', 180, 5, 6, 25, 2),
( 'hot_dog', 300, 12, 25, 2, 1),
( 'huevos_rancheros', 400, 20, 25, 30, 1),
( 'hummus', 250, 10, 12, 30, 0),
( 'ice_cream', 200, 3, 11, 28, 20),
( 'lasagna', 400, 25, 20, 40, 5),
( 'lobster_bisque', 300, 15, 20, 10, 1),
( 'lobster_roll_sandwich', 500, 30, 25, 40, 2),
( 'macaroni_and_cheese', 350, 12, 15, 45, 2),
( 'macarons', 200, 2, 8, 30, 20),
( 'miso_soup', 100, 5, 3, 12, 0),
( 'mussels', 200, 20, 6, 10, 0),
( 'nachos', 400, 10, 20, 50, 0),
( 'omelette', 250, 18, 18, 2, 1),
( 'onion_rings', 300, 4, 15, 40, 0),
( 'oysters', 50, 6, 2, 2, 0),
( 'pad_thai', 400, 20, 15, 55, 4),
( 'paella', 500, 25, 20, 60, 3),
( 'pancakes', 350, 8, 10, 60, 12),
( 'panna_cotta', 250, 5, 15, 28, 22),
( 'peking_duck', 600, 40, 40, 0, 0),
( 'pho', 350, 20, 10, 50, 1),
( 'pizza', 300, 12, 10, 36, 4),
( 'pork_chop', 400, 30, 20, 0, 0),
( 'poutine', 700, 25, 30, 90, 3),
( 'prime_rib', 600, 50, 40, 0, 0),
( 'pulled_pork_sandwich', 500, 30, 25, 45, 6),
( 'ramen', 400, 18, 15, 60, 4),
( 'ravioli', 300, 12, 10, 50, 3),
( 'red_velvet_cake', 350, 4, 20, 45, 28),
( 'risotto', 350, 10, 10, 60, 1),
( 'samosa', 250, 6, 12, 30, 3),
( 'sashimi', 200, 20, 5, 0, 0),
( 'scallops', 250, 30, 8, 0, 0),
( 'seaweed_salad', 150, 3, 5, 15, 0),
( 'shrimp_and_grits', 400, 20, 15, 50, 2),
( 'spaghetti_bolognese', 600, 25, 18, 80, 6),
( 'spaghetti_carbonara', 500, 20, 20, 60, 2),
( 'spring_rolls', 200, 4, 8, 30, 1),
( 'steak', 700, 50, 50, 0, 0),
( 'strawberry_shortcake', 300, 5, 10, 45, 18),
( 'sushi', 300, 15, 5, 45, 2),
( 'tacos', 250, 15, 10, 30, 1),
( 'takoyaki', 400, 12, 20, 45, 4),
( 'tiramisu', 450, 6, 24, 50, 25),
('tuna_tartare', 200, 22, 7, 0, 0),
('waffles', 350, 8, 14, 60, 14);

INSERT INTO menu_nutrition(menu_id, food_id) VALUES
(1, 6),  -- beet_salad
(2, 6),  
(4, 6),  
(1, 12), -- caesar_salad
(2, 12), 
(4, 12), 
(1, 14), -- caprese_salad
(2, 14), 
(4, 14), 
(1, 49), -- greek_salad
(2, 49), 
(4, 49), 
(1, 89), -- seaweed_salad
(2, 89), 
(4, 89),
(11, 96), -- sushi
(12, 96),  
(11, 87), -- sashimi
(12, 87),
(13, 76), -- Beef Bulgogi
(14, 76), -- Chicken Katsu
(15, 76), -- Vegetable Curry
(16, 76), -- Beef Pho
(10, 76); -- Pad Thai
