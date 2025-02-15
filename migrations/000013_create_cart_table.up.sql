-- Create the cart table
CREATE TABLE IF NOT EXISTS cart (
    id BIGSERIAL PRIMARY KEY,                            -- Unique identifier for the cart
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE, -- Reference to the user owning the cart
    total_price DOUBLE PRECISION NOT NULL DEFAULT 0,     -- Total price of items in the cart, default to 0
    status VARCHAR(20) DEFAULT 'active',                -- Status of the cart ('active', 'used', etc.)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,     -- Timestamp for when the cart was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP      -- Timestamp for when the cart was last updated
);

-- Create the cart_item table to link items or bundles to the cart
CREATE TABLE IF NOT EXISTS cart_item (
    id BIGSERIAL PRIMARY KEY,                            -- Unique identifier for the cart item
    cart_id BIGINT REFERENCES cart(id) ON DELETE CASCADE, -- Reference to the cart
    menu_id BIGINT NULL REFERENCES menu(id) ON DELETE CASCADE, -- Reference to the menu item (nullable if it's a bundling)
    bundling_id BIGINT NULL REFERENCES bundling(id) ON DELETE CASCADE, -- Reference to the bundling (nullable if it's a menu item)
    image_url VARCHAR(255),                              -- Image URL for the item/bundle
    name VARCHAR(225) NOT NULL,                          -- Name of the item/bundle
    quantity INT NOT NULL DEFAULT 1,                     -- Quantity
    price DOUBLE PRECISION NOT NULL,                     -- Price per unit of the item/bundle
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,      -- Timestamp for when the item was added
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,      -- Timestamp for when the item was last updated
    CHECK ((menu_id IS NOT NULL AND bundling_id IS NULL) OR (menu_id IS NULL AND bundling_id IS NOT NULL)) -- Ensure only one of menu_id or bundling_id is set
);

-- Trigger function to update `updated_at` timestamp for cart
CREATE OR REPLACE FUNCTION update_cart_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function to update `updated_at` timestamp for cart_item
CREATE OR REPLACE FUNCTION update_cart_item_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply triggers to update `updated_at` on row modifications
CREATE TRIGGER update_cart_updated_at
BEFORE UPDATE ON cart
FOR EACH ROW
EXECUTE FUNCTION update_cart_timestamp();

CREATE TRIGGER update_cart_item_updated_at
BEFORE UPDATE ON cart_item
FOR EACH ROW
EXECUTE FUNCTION update_cart_item_timestamp();

-- Seed data for cart table
INSERT INTO cart (user_id, total_price, status, created_at, updated_at)
VALUES
    (1, 350000, 'used', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),  -- Cart 1 for user_id 1
    (1, 150000, 'used', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);  -- Cart 2 for user_id 1

-- Insert multiple cart items for cart_id = 1 and cart_id = 2

INSERT INTO cart_item (cart_id, menu_id, bundling_id, image_url, name, quantity, price, created_at, updated_at)
VALUES
    -- cart_id = 1
    (1, 1, NULL, 'https://images.unsplash.com/photo-1670237735381-ac5c7fa72c51?q=80&w=2606&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 'Avocado Salad', 1, 45000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (1, 2, NULL, 'https://images.unsplash.com/photo-1604908176997-125f25cc6f3d?q=80&w=2626&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 'Grilled Chicken Salad', 1, 50000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (1, NULL, 1, 'https://images.unsplash.com/photo-1482049016688-2d3e1b311543?w=800&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTZ8fGZvb2R8ZW58MHx8MHx8fDA%3D', 'Langganan Mingguan', 1, 120000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (1, NULL, 2, 'https://images.unsplash.com/photo-1484980972926-edee96e0960d?w=800&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTh8fGZvb2R8ZW58MHx8MHx8fDA%3D', 'Langganan Bulanan', 1, 500000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

    -- cart_id = 2
    (2, 3, NULL, 'https://images.unsplash.com/photo-1588166524941-3bf61a9c41db?q=80&w=2784&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 'Spaghetti Carbonara', 2, 70000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, 4, NULL, 'https://plus.unsplash.com/premium_photo-1700089483464-4f76cc3d360b?q=80&w=2787&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 'Chicken Caesar Salad', 1, 60000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, NULL, 1, 'https://images.unsplash.com/photo-1482049016688-2d3e1b311543?w=800&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTZ8fGZvb2R8ZW58MHx8MHx8fDA%3D', 'Langganan Mingguan', 1, 120000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, NULL, 2, 'https://images.unsplash.com/photo-1484980972926-edee96e0960d?w=800&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTh8fGZvb2R8ZW58MHx8MHx8fDA%3D', 'Langganan Bulanan', 1, 500000, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
