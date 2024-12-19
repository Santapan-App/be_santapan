-- Table for couriers
CREATE TABLE IF NOT EXISTS couriers (
    id BIGSERIAL PRIMARY KEY,              -- Unique identifier for the courier
    name VARCHAR(255) NOT NULL,            -- Name of the courier (e.g., 'DHL', 'FedEx')
    logo VARCHAR(255) NOT NULL,            -- Logo URL or path for the courier's logo
    price DECIMAL(10, 2) NOT NULL,         -- Price for the courier service
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Timestamp when the courier was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- Timestamp for the last update
);

-- Trigger function to update `updated_at` timestamp for couriers
CREATE OR REPLACE FUNCTION update_courier_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();  -- Update `updated_at` to current time
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to update `updated_at` on row modifications
CREATE TRIGGER update_courier_updated_at
BEFORE UPDATE ON couriers
FOR EACH ROW
EXECUTE FUNCTION update_courier_timestamp();

-- Insert seed data for couriers
INSERT INTO couriers (name, logo, price) VALUES
    ('GOJEK', 'https://blogger.googleusercontent.com/img/b/R29vZ2xl/AVvXsEhtKkngF0s65WHeRxAwldlEtDGY0nbh1OZQSq21A3RiMv24rFOdi05MjFLxbmK2iw39jcp7ElLlkSw6MuMyedhWTHM883wxSapK6HwEqJeyhX4wvCgofb_YlROPAQF2xCIBUFFR7AMPmIrg/s2048/Gojek.png', 10000),
    ('GRAB', 'https://www.carnivalworld.sg/wp-content/uploads/2021/10/Grab-logo.png', 15000);