-- Insert seed data into the users table
INSERT INTO users (full_name, email, password, created_at, updated_at, email_verified_at)
VALUES
('Jan Falih Fadhillah', 'bosspulsa57@gmail.com', '$2a$12$K.vnPeV0Hlm/MDJcV5lV9.Txfb8Bniw6blg7BdM2aIfIghA4Jc.OK', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
('Pahala', 'pfnazhmi@gmail.com', '$2a$12$K.vnPeV0Hlm/MDJcV5lV9.Txfb8Bniw6blg7BdM2aIfIghA4Jc.OK', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Robert Brown', 'robert.brown@example.com', 'hashed_password_3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
