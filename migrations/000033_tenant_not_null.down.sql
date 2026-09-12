ALTER TABLE user_preferences DROP CONSTRAINT IF EXISTS user_preferences_pkey;
CREATE UNIQUE INDEX user_preferences_tenant_user_scope_key
  ON user_preferences(tenant_id, user_id, scope) WHERE tenant_id IS NOT NULL;

DROP INDEX IF EXISTS courses_tenant_slug_key;
CREATE UNIQUE INDEX courses_tenant_slug_key
  ON courses(tenant_id, slug) WHERE tenant_id IS NOT NULL;

ALTER TABLE user_preferences ALTER COLUMN tenant_id DROP NOT NULL;
ALTER TABLE notifications ALTER COLUMN tenant_id DROP NOT NULL;
ALTER TABLE conversations ALTER COLUMN tenant_id DROP NOT NULL;
ALTER TABLE calendar_events ALTER COLUMN tenant_id DROP NOT NULL;
ALTER TABLE courses ALTER COLUMN tenant_id DROP NOT NULL;
