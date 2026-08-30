ALTER TABLE campus_profiles
  ADD COLUMN IF NOT EXISTS is_blocked BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_campus_profiles_search
  ON campus_profiles (profile_type, full_name);
