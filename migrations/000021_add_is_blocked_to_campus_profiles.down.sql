DROP INDEX IF EXISTS idx_campus_profiles_search;
ALTER TABLE campus_profiles DROP COLUMN IF EXISTS is_blocked;
