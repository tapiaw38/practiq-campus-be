CREATE TABLE course_groups (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  name VARCHAR(120) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_course_groups_course ON course_groups(course_id);

CREATE TABLE course_group_members (
  group_id UUID NOT NULL REFERENCES course_groups(id) ON DELETE CASCADE,
  user_id VARCHAR(255) NOT NULL,
  PRIMARY KEY (group_id, user_id)
);
