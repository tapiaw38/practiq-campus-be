CREATE TABLE user_preferences (
  user_id VARCHAR(255) NOT NULL REFERENCES campus_profiles(id) ON DELETE CASCADE,
  scope TEXT NOT NULL,
  settings JSONB NOT NULL DEFAULT '{}'::jsonb,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, scope)
);

CREATE INDEX idx_user_preferences_scope ON user_preferences(scope);
