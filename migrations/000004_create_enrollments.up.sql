-- enrollment_role is course-scoped, not a JWT-level role: it grants elevated
-- access (teaching_assistant/co_teacher) within THIS course only, alongside
-- the user's global role, so a student in one course can TA another without
-- a new global role or a full course-level permissions table.
CREATE TABLE IF NOT EXISTS enrollments (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id       UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    user_id         VARCHAR(255) NOT NULL REFERENCES campus_profiles(id) ON DELETE CASCADE,
    enrollment_role VARCHAR(30) NOT NULL DEFAULT 'student' CHECK (enrollment_role IN ('student', 'teaching_assistant', 'co_teacher')),
    status          VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'dropped', 'completed')),
    enrolled_at     TIMESTAMP DEFAULT NOW(),
    UNIQUE(course_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_enrollments_course_id ON enrollments(course_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_user_id ON enrollments(user_id);
