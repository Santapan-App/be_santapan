-- Table for transactions
CREATE TABLE IF NOT EXISTS transaction (
    id BIGSERIAL PRIMARY KEY,                       -- Unique identifier for the transaction
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,                        -- Reference to the user who made the transaction
    cart_id BIGINT REFERENCES cart(id) ON DELETE CASCADE,                         -- Reference to the cart used in the transaction
    payment_id BIGINT REFERENCES payment(id) ON DELETE CASCADE,                     -- Reference to the payment transaction
    courier_id BIGINT REFERENCES couriers(id) ON DELETE CASCADE,                     -- Reference to the courier used for the transaction
    address_id BIGINT REFERENCES address(id) ON DELETE CASCADE,                     -- Reference to the address used for the transaction
    amount DECIMAL(10, 2) NOT NULL,            -- Total price of the transaction
    status VARCHAR(50) NOT NULL,        -- e.g., 'unpaid', 'ongoing', 'completed', 'failed'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp for when the transaction was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Timestamp for the last update
);

-- Trigger function to update `updated_at` timestamp for transaction
CREATE OR REPLACE FUNCTION update_transaction_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Apply triggers to update `updated_at` on row modifications
CREATE TRIGGER update_transaction_updated_at
BEFORE UPDATE ON transaction
FOR EACH ROW
EXECUTE FUNCTION update_transaction_timestamp();


-- Seed data for transaction table
INSERT INTO transaction (user_id, cart_id, payment_id, courier_id, address_id, amount, status, created_at, updated_at)
VALUES
    (1, 1, 1, 1, 1, 350000, 'failed', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),  -- Transaction 1 for user_id 1, cart_id 1
    (1, 2, 1, 1, 1, 150000, 'success', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);  -- Transaction 2 for user_id 1, cart_id 2
