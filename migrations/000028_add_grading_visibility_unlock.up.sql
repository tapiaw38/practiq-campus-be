ALTER TABLE assignments ADD COLUMN weight INT NOT NULL DEFAULT 100 CHECK (weight > 0);
ALTER TABLE assignments ADD COLUMN visible_group_id UUID REFERENCES course_groups(id) ON DELETE SET NULL;
ALTER TABLE assignments ADD COLUMN unlock_after_type VARCHAR(20);
ALTER TABLE assignments ADD COLUMN unlock_after_id UUID;

ALTER TABLE quizzes ADD COLUMN weight INT NOT NULL DEFAULT 100 CHECK (weight > 0);
ALTER TABLE quizzes ADD COLUMN visible_group_id UUID REFERENCES course_groups(id) ON DELETE SET NULL;
ALTER TABLE quizzes ADD COLUMN unlock_after_type VARCHAR(20);
ALTER TABLE quizzes ADD COLUMN unlock_after_id UUID;

CREATE INDEX idx_assignments_visible_group ON assignments(visible_group_id) WHERE visible_group_id IS NOT NULL;
CREATE INDEX idx_quizzes_visible_group ON quizzes(visible_group_id) WHERE visible_group_id IS NOT NULL;
