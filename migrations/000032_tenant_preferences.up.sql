-- Preferences are per person per institution. Somebody who teaches at two of
-- them keeps a separate setup in each: the same account, two contexts.
--
-- No backfill: Campus has no data yet, and a preference with no tenant would
-- be a row no context can claim.
ALTER TABLE user_preferences ADD COLUMN tenant_id UUID REFERENCES campus_tenants(id);

-- The primary key made a preference global to the account, so selecting
-- another institution would have overwritten the first one's settings.
--
-- It is replaced by a partial unique index rather than a new primary key
-- because tenant_id is still nullable here; the migration that promotes it to
-- NOT NULL restores a proper key over all three columns.
ALTER TABLE user_preferences DROP CONSTRAINT IF EXISTS user_preferences_pkey;
CREATE UNIQUE INDEX user_preferences_tenant_user_scope_key
  ON user_preferences(tenant_id, user_id, scope) WHERE tenant_id IS NOT NULL;
CREATE INDEX idx_user_preferences_tenant_user ON user_preferences(tenant_id, user_id);
