CREATE TABLE IF NOT EXISTS courses (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id    VARCHAR(255) NOT NULL REFERENCES campus_profiles(id) ON DELETE CASCADE,
    title       VARCHAR(200) NOT NULL,
    slug        VARCHAR(220) NOT NULL UNIQUE,
    description TEXT DEFAULT '',
    status      VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    start_date  DATE,
    end_date    DATE,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_courses_owner_id ON courses(owner_id);
CREATE INDEX IF NOT EXISTS idx_courses_status ON courses(status);
