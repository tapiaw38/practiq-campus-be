DROP INDEX IF EXISTS idx_course_materials_assignment;
ALTER TABLE course_materials DROP COLUMN IF EXISTS assignment_id;
