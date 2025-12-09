-- Seed the initial admin user
-- Password hash should be generated and inserted separately via the app
-- This just creates the user record

INSERT INTO users (email, name, role)
VALUES ('kellon08@gmail.com', 'Kellon Sandall', 'admin')
ON CONFLICT (email) DO NOTHING;

-- Note: The login record with password_hash must be created by the application
-- using bcrypt. Run the seed command after migrations to set the initial password.
