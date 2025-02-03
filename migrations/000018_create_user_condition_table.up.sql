-- Table for user_condition
CREATE TABLE IF NOT EXISTS user_condition (
    id BIGSERIAL PRIMARY KEY,                       -- Unique identifier for the transaction
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE, -- Reference to the user who made the transaction
    diabetes BOOLEAN DEFAULT FALSE,                 -- Indicates if diabetes is relevant to the transaction
    gerd BOOLEAN DEFAULT FALSE,                     -- Indicates if GERD is relevant to the transaction
    asam_urat BOOLEAN DEFAULT FALSE,                -- Indicates if asam urat (gout) is relevant to the transaction
    kolestrol BOOLEAN DEFAULT FALSE,                -- Indicates if kolesterol is relevant to the transaction
    rendah_karbohidrat BOOLEAN DEFAULT FALSE,       -- Indicates if a low carbohydrate preference is applied
    tinggi_protein BOOLEAN DEFAULT FALSE,           -- Indicates if a high protein preference is applied
    vegetarian BOOLEAN DEFAULT FALSE,               -- Indicates if a vegetarian preference is applied
    rendah_gula BOOLEAN DEFAULT FALSE,              -- Indicates if a low sugar preference is applied
    rendah_kalori BOOLEAN DEFAULT FALSE,            -- Indicates if a low calorie preference is applied
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp for when the transaction was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Timestamp for the last update
    CONSTRAINT unique_user_id UNIQUE (user_id) -- Ensure each user has only one record in this table
);

-- Trigger function to update `updated_at` timestamp for transaction
CREATE OR REPLACE FUNCTION update_transaction_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;

$$ LANGUAGE plpgsql;

-- Trigger to call the update function on any update of the user_condition table
CREATE TRIGGER update_user_condition_timestamp
BEFORE UPDATE ON user_condition
FOR EACH ROW
EXECUTE FUNCTION update_transaction_timestamp();
