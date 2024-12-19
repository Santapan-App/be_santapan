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
