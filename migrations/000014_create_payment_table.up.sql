-- Table for payments
CREATE TABLE IF NOT EXISTS payment (
    id BIGSERIAL PRIMARY KEY,            -- Unique identifier for the payment method
    reference_id VARCHAR(50) NOT NULL,   -- Reference ID for the payment method (e.g., 'CC', 'BT', 'EW')
    session_id VARCHAR(255) NOT NULL,     -- Session ID for the payment method
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE, -- Reference to the user who owns the payment method
    amount DECIMAL(10, 2) NOT NULL,       -- Total amount of the payment
    status VARCHAR(50) NOT NULL,          -- Status of the payment (e.g., 'unpaid', 'ongoing', 'completed', 'failed')
    url VARCHAR(255) NOT NULL,            -- URL for the payment method
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp for when the payment method was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Timestamp for the last update
);

-- Trigger function to update `updated_at` timestamp for payment
CREATE OR REPLACE FUNCTION update_payment_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to update `updated_at` on row modifications
CREATE TRIGGER update_payment_updated_at
BEFORE UPDATE ON payment
FOR EACH ROW
EXECUTE FUNCTION update_payment_timestamp();
