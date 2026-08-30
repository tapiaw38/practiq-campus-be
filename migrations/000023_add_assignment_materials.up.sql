ALTER TABLE course_materials ADD COLUMN assignment_id UUID REFERENCES assignments(id) ON DELETE CASCADE;
CREATE INDEX idx_course_materials_assignment ON course_materials(assignment_id, created_at DESC);
