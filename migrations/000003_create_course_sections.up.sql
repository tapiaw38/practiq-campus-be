CREATE TABLE IF NOT EXISTS course_sections (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id  UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title      VARCHAR(200) NOT NULL,
    position   INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_course_sections_course_id ON course_sections(course_id);
