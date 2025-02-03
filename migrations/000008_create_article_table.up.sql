-- Create the article table
CREATE TABLE IF NOT EXISTS article (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,                  -- Maps to `Title` field in Go
    content TEXT NOT NULL,                        -- Maps to `Content` field in Go
    image_url VARCHAR(500),                       -- Maps to `ImageURL` field in Go
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create an index for the id column for faster lookups
CREATE INDEX idx_article_id ON article(id);

-- Create or replace the function to update the updated_at column
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Drop the existing trigger if it exists
DROP TRIGGER IF EXISTS update_article_updated_at ON article;

-- Create the trigger to call the function before update
CREATE TRIGGER update_article_updated_at
BEFORE UPDATE ON article
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Insert sample data into the article table with content about the Santapan app
INSERT INTO article (title, content, image_url, created_at, updated_at)
VALUES 
    ('Welcome to Santapan: Your Guide to Healthy Eating', 
     'Santapan is designed to help users make informed choices about their diets. Discover how this app can support your journey toward balanced and nutritious meals.', 
     'https://plus.unsplash.com/premium_photo-1661777702966-aed29ab4106b?q=80&w=2940&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 
     CURRENT_TIMESTAMP, 
     CURRENT_TIMESTAMP),
     
    ('5 Reasons to Use Santapan for Meal Planning', 
     'Meal planning is crucial for maintaining a healthy diet. Santapan offers customizable meal plans, nutrition tracking, and recipe suggestions to simplify the process.', 
     'https://plus.unsplash.com/premium_photo-1723291306365-841e10185b6f?w=800&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NXx8bWVhbHxlbnwwfHwwfHx8MA%3D%3D', 
     CURRENT_TIMESTAMP, 
     CURRENT_TIMESTAMP),
     
    ('Understanding Nutritional Labels with Santapan', 
     'Deciphering nutritional labels can be challenging. Santapan provides an easy guide to understand what’s in your food, helping you make healthier choices.', 
     'https://images.unsplash.com/photo-1628191138144-a51eeee8e2c3?w=800&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTB8fG1lYWx8ZW58MHx8MHx8fDA%3D', 
     CURRENT_TIMESTAMP, 
     CURRENT_TIMESTAMP),

    ('How Santapan Supports Local Farmers and Fresh Ingredients', 
     'Santapan partners with local farmers to promote fresh and organic ingredients. Learn how the app connects users with nearby sources for healthier options.', 
     'https://plus.unsplash.com/premium_photo-1724129050515-2d756cdeeb72?w=800&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8OXx8bG9jYWwlMjBmYXJtZXJzfGVufDB8fDB8fHww', 
     CURRENT_TIMESTAMP, 
     CURRENT_TIMESTAMP),

    ('Personalized Nutrition Insights with Santapan', 
     'Everyone’s nutritional needs are different. Santapan’s personalized insights help you understand the best dietary choices based on your health goals and preferences.', 
     'https://plus.unsplash.com/premium_photo-1661431006380-ef37401c681b?w=800&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MXx8bWVhbCUyMGluc2lnaHR8ZW58MHx8MHx8fDA%3D', 
     CURRENT_TIMESTAMP, 
     CURRENT_TIMESTAMP);
