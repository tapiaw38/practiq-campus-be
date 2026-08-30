CREATE TABLE IF NOT EXISTS forum_threads (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id  UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    author_id  VARCHAR(255) NOT NULL REFERENCES campus_profiles(id) ON DELETE CASCADE,
    title      VARCHAR(200) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_forum_threads_course_id ON forum_threads(course_id);
