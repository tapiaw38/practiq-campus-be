CREATE TABLE notifications (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), user_id UUID NOT NULL, type VARCHAR(40) NOT NULL, title VARCHAR(255) NOT NULL, body TEXT NOT NULL DEFAULT '', data JSONB NOT NULL DEFAULT '{}'::jsonb, read_at TIMESTAMPTZ, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE INDEX idx_notifications_user_created ON notifications(user_id, created_at DESC);
