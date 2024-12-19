CREATE TABLE IF NOT EXISTS address (
    id BIGSERIAL PRIMARY KEY,               -- Unique identifier for the address
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,  -- Reference to the user who owns the address
    label VARCHAR(100) NOT NULL,            -- Label for the address (e.g., 'Home', 'Work')
    address TEXT NOT NULL,                  -- The actual address
    name VARCHAR(255) NOT NULL,             -- Name of the person associated with the address
    notes TEXT,                             -- Any additional notes about the address
    phone VARCHAR(20) NOT NULL,             -- Phone number associated with the address
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Timestamp when the address was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- Timestamp when the address was last updated
);

-- Trigger function to update `updated_at` timestamp for address
CREATE OR REPLACE FUNCTION update_address_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;

$$ LANGUAGE plpgsql;

-- Apply trigger to update `updated_at` on row modifications
CREATE TRIGGER update_address_updated_at
BEFORE UPDATE ON address
FOR EACH ROW
EXECUTE FUNCTION update_address_timestamp();
