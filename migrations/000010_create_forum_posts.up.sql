CREATE TABLE IF NOT EXISTS forum_posts (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    thread_id  UUID NOT NULL REFERENCES forum_threads(id) ON DELETE CASCADE,
    author_id  VARCHAR(255) NOT NULL REFERENCES campus_profiles(id) ON DELETE CASCADE,
    body       TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_forum_posts_thread_id ON forum_posts(thread_id);
