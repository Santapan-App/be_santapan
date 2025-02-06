-- Create ENUM type for bundling_type
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bundling_type_enum') THEN
        CREATE TYPE bundling_type_enum AS ENUM ('weekly', 'monthly');
    END IF;
END;
$$;

-- Table for menu items
CREATE TABLE IF NOT EXISTS menu (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    image_url VARCHAR(255),
    nutrition JSONB,       -- JSON format for flexible nutrition details
    features JSONB,        -- Features stored as JSON for flexibility
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Table for link category and menu items
CREATE TABLE IF NOT EXISTS category_menu (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT REFERENCES category(id) ON DELETE CASCADE,
    menu_id BIGINT REFERENCES menu(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(category_id, menu_id)  -- Ensure unique combination of category and menu items
);

-- Table for bundling types
CREATE TABLE IF NOT EXISTS bundling (
    id BIGSERIAL PRIMARY KEY,
    bundling_name VARCHAR(255),
    bundling_type bundling_type_enum NOT NULL,  -- ENUM type for 'weekly' or 'monthly'
    price DECIMAL(10, 2) NOT NULL,
    image_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table to link bundling and menu items
CREATE TABLE IF NOT EXISTS bundling_menu (
    id BIGSERIAL PRIMARY KEY,
    bundling_id BIGINT REFERENCES bundling(id) ON DELETE CASCADE,
    menu_id BIGINT REFERENCES menu(id) ON DELETE CASCADE,
    day_number INT NOT NULL,          -- Represents the day number (1-7 for weekly, 1-30 for monthly)
    meal_description TEXT,            -- Description of the meal for the day
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(bundling_id, menu_id, day_number)  -- Ensure unique combination of bundling and menu items per day
);

-- Trigger function to update `updated_at` timestamp
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to update `updated_at` column on row modifications
CREATE TRIGGER update_menu_timestamp
BEFORE UPDATE ON menu
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_bundling_timestamp
BEFORE UPDATE ON bundling
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_bundling_menu_timestamp
BEFORE UPDATE ON bundling_menu
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Seed data for menu table
INSERT INTO menu (title, description, price, image_url, nutrition, features)
VALUES
    ('Avocado Salad', 'A healthy avocado salad with fresh vegetables and a side of toast', 45000, 'https://www.recipetineats.com/uploads/2023/11/Avocado-salad_0.jpg', '{"calories": 350, "protein": 5, "fat": 30}', '{"gluten_free": true, "vegetarian": true}'),
    ('Grilled Chicken Salad', 'A delicious grilled chicken salad', 50000, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTlbShllQYeV7L3R5hRsnpluh-9g5oQGssnMA&s', '{"calories": 400, "protein": 35, "fat": 0}', '{"high_protein": true}'),
    ('Spaghetti Carbonara', 'A creamy pasta dish with bacon and parmesan cheese', 70000, 'https://asset.kompas.com/crops/534CByKLuUtN81KtzF3RZkkil8A=/249x82:1000x583/1200x800/data/photo/2022/12/27/63aaef8f6a2f6.jpeg', '{"calories": 850, "protein": 20, "fat": 35}', '{"gluten_free": false}'),
    ('Chicken Caesar Salad', 'Grilled chicken with Caesar dressing and croutons', 60000, 'https://s23209.pcdn.co/wp-content/uploads/2023/01/220905_DD_Chx-Caesar-Salad_051-500x500.jpg', '{"calories": 450, "protein": 30, "fat": 20}', '{"gluten_free": false}'),
    ('Tofu Stir-fry', 'A healthy stir-fry with tofu, vegetables, and soy sauce', 50000, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQVJ-l7PY4VDYMV-65pgIjEuCPOKvtzzo1J_Q&s', '{"calories": 400, "protein": 15, "fat": 18}', '{"vegetarian": true, "vegan": true}'),
    ('Chocolate Cake', 'A rich and moist chocolate cake with frosting', 30000, 'https://thebigmansworld.com/wp-content/uploads/2020/08/3-layer-chocolate-cake7734.jpg', '{"calories": 350, "protein": 4, "fat": 18}', '{"vegetarian": true}'),
    ('Fried Rice', 'Classic fried rice with vegetables and a choice of meat', 55000, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSBNVqnMHO7y8bpvhmCBFfDwOesFr7PR6LXkA&s', '{"calories": 600, "protein": 20, "fat": 22}', '{"gluten_free": false}'),
    ('Miso Soup', 'Traditional Japanese miso soup with tofu and seaweed', 20000, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQh4xEwOLtDhKjT_c0hARFZcBQ9y-v2Nrv-Dw&s', '{"calories": 150, "protein": 8, "fat": 5}', '{"gluten_free": true, "vegetarian": true}'),
    ('Tempura', 'Crispy tempura shrimp with a dipping sauce', 70000, 'https://media.istockphoto.com/id/184274745/id/foto/tempura-udang.jpg?s=612x612&w=0&k=20&c=LsLUFmlegF7CgJiHfKD0IJ102SaUofbsUDCbjScCZfE=', '{"calories": 500, "protein": 15, "fat": 25}', '{"gluten_free": false}'),
    ('Pad Thai', 'Traditional Thai stir-fried rice noodles with shrimp, tofu, and peanuts', 55000, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTvvnHh_LTnQFZIil8QjwC7dYCSDiHXiV3Z7g&s', '{"calories": 400, "protein": 20, "fat": 15, "carbohydrates": 55, "sugar": 4}', '{"gluten_free": false, "high_protein": true}'),
    ('Grilled Salmon', 'Perfectly grilled salmon fillet served with a side of vegetables', 75000, 'https://www.thecookierookie.com/wp-content/uploads/2023/05/featured-grilled-salmon-recipe.jpg', '{"calories": 350, "protein": 40, "fat": 20, "carbohydrates": 0, "sugar": 0}', '{"gluten_free": true, "high_protein": true}'),
    ('Sushi Platter', 'Assorted sushi rolls and sashimi with soy sauce and wasabi', 80000, 'https://images.unsplash.com/photo-1604908176997-125f25cc6f3d?q=80&w=2626&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', '{"calories": 500, "protein": 25, "fat": 15, "carbohydrates": 60, "sugar": 5}', '{"gluten_free": false}'),
    ('Beef Bulgogi', 'Korean-style marinated beef with rice and kimchi', 65000, 'https://asset.kompas.com/crops/0P7PRcCKawSEtEbla10eqVHeiqE=/0x0:0x0/1200x800/data/photo/2020/12/23/5fe2c9413a6ce.jpg', '{"calories": 450, "protein": 30, "fat": 20, "carbohydrates": 40, "sugar": 10}', '{"gluten_free": false, "high_protein": true}'),
    ('Chicken Katsu', 'Crispy breaded chicken cutlet with tonkatsu sauce', 60000, 'https://oishipedia.com/wp-content/uploads/2023/09/Snapinsta.app_199642976_1442201386143212_7434957661754949293_n_1080-2.jpg', '{"calories": 400, "protein": 25, "fat": 15, "carbohydrates": 30, "sugar": 5}', '{"gluten_free": false, "high_protein": true}'),
    ('Vegetable Curry', 'Spicy vegetable curry with rice and naan bread', 50000, 'https://images.immediate.co.uk/production/volatile/sites/30/2022/06/Courgette-curry-c295fa0.jpg?resize=768,574', '{"calories": 350, "protein": 10, "fat": 15, "carbohydrates": 40, "sugar": 5}', '{"gluten_free": false, "vegetarian": true}'),
    ('Beef Pho', 'Vietnamese beef noodle soup with fresh herbs and lime', 55000, 'https://www.recipetineats.com/tachyon/2019/04/Beef-Pho_6.jpg', '{"calories": 400, "protein": 20, "fat": 10, "carbohydrates": 50, "sugar": 5}', '{"gluten_free": false, "high_protein": true}'),
    ('Falafel Wrap', 'Middle Eastern falafel with hummus and tahini in a wrap', 45000, 'https://static.toiimg.com/thumb/62708678.cms?imgsize=156976&width=800&height=800', '{"calories": 350, "protein": 10, "fat": 15, "carbohydrates": 40, "sugar": 5}', '{"gluten_free": false, "vegetarian": true}');

-- Seed data for category_menu table to link menu items to categories
-- Assuming category_id of 1 corresponds to the 'Indonesian' category in the `category` table
-- Assuming menu_id of 1 corresponds to the 'Avocado Salad' menu item in the `menu` table
-- Seed data for category_menu table to link menu items to categories
INSERT INTO category_menu (category_id, menu_id)
VALUES
    -- Indonesian (category_id = 1)
    (1, 1),  -- Avocado Salad
    (1, 2),  -- Grilled Chicken Salad
    (1, 3),  -- Spaghetti Carbonara
    (1, 4),  -- Chicken Caesar Salad

    -- Western (category_id = 2)
    (2, 5),  -- Tofu Stir-fry
    (2, 6),  -- Chocolate Cake
    (2, 7),  -- Fried Rice

    -- Chinese (category_id = 3)
    (3, 8),  -- Miso Soup
    (3, 9),  -- Tempura

    -- Japanese (category_id = 4)
    (4, 11), -- Grilled Salmon
    (4, 12), -- Sushi Platter
    (4, 13), -- Beef Bulgogi

    -- Desserts (category_id = 5)
    (5, 6),  -- Chocolate Cake
    (5, 14), -- Chicken Katsu
    (5, 15), -- Vegetable Curry

    -- Vegan (category_id = 6)
    (6, 5),  -- Tofu Stir-fry
    (6, 15), -- Vegetable Curry
    (6, 1),
    (7, 11), -- Grilled Salmon
    (7, 12), -- Sushi Platter
    (7, 17); -- Falafel Wrap
    
-- Seed data for bundling table
INSERT INTO bundling (bundling_name, bundling_type, price, image_url)
VALUES
    ('Langganan Mingguan', 'weekly', 120000, 'https://i.ibb.co.com/4Z1353rg/langganan-migguan.png'),
    ('Langganan Bulanan', 'monthly', 500000, 'https://i.ibb.co.com/FLK72b8Z/langganan-bulanan.png');

-- Seed data for bundling_menu table to represent meals for each day of the week
-- Assuming bundling_id of 1 corresponds to the 'weekly' bundling in the `bundling` table
INSERT INTO bundling_menu (bundling_id, menu_id, day_number, meal_description)
VALUES
    (1, 1, 1, 'Lunch - Avocado Salad'),
    (1, 2, 1, 'Lunch - Grilled Chicken Salad'),
    (1, 1, 2, 'Lunch - Avocado Salad'),
    (1, 1, 3, 'Lunch - Avocado Salad'),
    (1, 2, 4, 'Lunch - Grilled Chicken Salad'),
    (1, 2, 5, 'Lunch - Grilled Chicken Salad'),
    (1, 2, 6, 'Lunch - Grilled Chicken Salad'),
    (1, 1, 7, 'Lunch - Avocado Salad'),
    --> Monthly
    (2, 1, 1, 'Salad with Avocado and Fresh Vegetables'),
    (2, 2, 1, 'Chicken with Fresh Vegetables and Rice'),
    (2, 1, 2, 'Fresh Avocado Salad with Crunchy Toast'),
    (2, 1, 3, 'Healthy Avocado Salad with Seasonal Greens'),
    (2, 2, 4, 'Grilled Chicken with Mixed Vegetables'),
    (2, 2, 5, 'Classic Grilled Chicken with Rice and Vegetables'),
    (2, 2, 6, 'Succulent Grilled Chicken with Herbed Rice'),
    (2, 1, 7, 'Avocado Salad with Olive Oil Dressing'),
    (2, 1, 8, 'Refreshing Avocado Salad with Lemon Vinaigrette'),
    (2, 2, 8, 'Tender Grilled Chicken with Sweet Corn'),
    (2, 3, 9, 'Spaghetti Carbonara with Creamy Parmesan Sauce'),
    (2, 4, 10, 'Chicken Caesar Salad with Crispy Croutons'),
    (2, 5, 11, 'Stir-fried Tofu with Bell Peppers and Soy Sauce'),
    (2, 6, 12, 'Rich Chocolate Cake with Ganache Frosting'),
    (2, 7, 13, 'Fried Rice with Egg and Mixed Vegetables'),
    (2, 1, 14, 'Avocado Salad with Cherry Tomatoes and Cucumbers'),
    (2, 2, 15, 'Grilled Chicken Breast with Steamed Broccoli'),
    (2, 3, 16, 'Creamy Spaghetti Carbonara with Bacon'),
    (2, 4, 17, 'Chicken Caesar Salad with Parmesan Shavings'),
    (2, 5, 18, 'Vegan Tofu Stir-fry with Spicy Garlic Sauce'),
    (2, 6, 19, 'Decadent Chocolate Cake with Whipped Cream'),
    (2, 7, 20, 'Classic Fried Rice with Shrimp and Peas'),
    (2, 1, 21, 'Light Avocado Salad with Spinach and Walnuts'),
    (2, 2, 22, 'Grilled Chicken Thighs with Honey Glaze'),
    (2, 3, 23, 'Spaghetti Carbonara with Truffle Oil'),
    (2, 4, 24, 'Caesar Salad with Grilled Chicken Strips'),
    (2, 5, 25, 'Tofu Stir-fry with Carrots and Snow Peas'),
    (2, 6, 26, 'Moist Chocolate Cake with Strawberry Garnish'),
    (2, 7, 27, 'Fried Rice with Pineapple and Cashews'),
    (2, 1, 28, 'Avocado Salad with Feta Cheese and Basil'),
    (2, 2, 29, 'Grilled Chicken with Garlic Butter Sauce'),
    (2, 3, 30, 'Spaghetti Carbonara with Fresh Parsley');
