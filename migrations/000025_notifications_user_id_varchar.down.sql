ALTER TABLE notifications ALTER COLUMN user_id TYPE UUID USING user_id::uuid;
