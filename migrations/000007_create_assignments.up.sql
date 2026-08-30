CREATE TABLE IF NOT EXISTS assignments (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id   UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    section_id  UUID REFERENCES course_sections(id) ON DELETE SET NULL,
    title       VARCHAR(200) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    due_at      TIMESTAMP,
    max_score   INT NOT NULL DEFAULT 100,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_assignments_course_id ON assignments(course_id);
