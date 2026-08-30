CREATE TABLE IF NOT EXISTS calendar_events (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id   VARCHAR(255) NOT NULL REFERENCES campus_profiles(id) ON DELETE CASCADE,
    course_id  UUID REFERENCES courses(id) ON DELETE CASCADE,
    title      VARCHAR(200) NOT NULL,
    starts_at  TIMESTAMP NOT NULL,
    ends_at    TIMESTAMP,
    all_day    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_calendar_events_owner_id ON calendar_events(owner_id);
