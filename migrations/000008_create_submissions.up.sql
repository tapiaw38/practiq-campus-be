CREATE TABLE IF NOT EXISTS submissions (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assignment_id UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    user_id       VARCHAR(255) NOT NULL REFERENCES campus_profiles(id) ON DELETE CASCADE,
    content       TEXT NOT NULL DEFAULT '',
    status        VARCHAR(20) NOT NULL DEFAULT 'submitted' CHECK (status IN ('submitted', 'graded')),
    score         INT,
    feedback      TEXT NOT NULL DEFAULT '',
    submitted_at  TIMESTAMP DEFAULT NOW(),
    graded_at     TIMESTAMP,
    UNIQUE (assignment_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_submissions_assignment_id ON submissions(assignment_id);
