-- Create the banner table
CREATE TABLE IF NOT EXISTS category (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,               -- Maps to `Title` field in Go
    image_url VARCHAR(255),                     -- Maps to `ImageURL` field in Go
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create an index for the id column for faster lookups
CREATE INDEX idx_category_id ON category(id);

-- Create or replace the function to update the updated_at column
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Drop the existing trigger if it exists
DROP TRIGGER IF EXISTS update_category_updated_at ON category;

-- Create the trigger to call the function before update
CREATE TRIGGER update_category_updated_at
BEFORE UPDATE ON category
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();


-- Insert sample data into the category table for food categories in Santapan
INSERT INTO category (title, image_url)
VALUES
    ('Indonesian', 'https://images.unsplash.com/photo-1564671165093-20688ff1fffa?q=80&w=2766&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'),
    ('Western', 'https://images.unsplash.com/photo-1686994562893-001a18e5bc43?q=80&w=2960&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'),
    ('Chinese', 'https://plus.unsplash.com/premium_photo-1661600135596-dcb910b05340?q=80&w=2942&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'),
    ('Japanese', 'https://plus.unsplash.com/premium_photo-1723874570807-570c56b41e4e?q=80&w=2940&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'),
    ('Desserts', 'https://images.unsplash.com/photo-1495147466023-ac5c588e2e94?q=80&w=2787&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'),
    ('Vegan', 'https://images.unsplash.com/photo-1623428187969-5da2dcea5ebf?q=80&w=2864&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'),
    ('Seafood', 'https://plus.unsplash.com/premium_photo-1707935175109-ba307d98bfe2?q=80&w=2787&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'),
    ('Asian', 'https://i.pinimg.com/736x/5e/af/12/5eaf1298fe096154a3771ad40ace04b1.jpg');;