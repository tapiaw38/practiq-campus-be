-- Every write path is tenantized now, so a row with no tenant can only be one
-- nothing created: the constraint states what the code already guarantees, and
-- turns a future omission into an error at the database rather than a row no
-- institution can claim.
--
-- No backfill. Campus has no production data, and inventing a tenant for a row
-- would be guessing which institution owns it.
ALTER TABLE courses ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE calendar_events ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE conversations ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE notifications ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE user_preferences ALTER COLUMN tenant_id SET NOT NULL;

-- With the column guaranteed, the partial indexes can become plain ones and
-- the preference key can go back to being a primary key.
DROP INDEX IF EXISTS courses_tenant_slug_key;
CREATE UNIQUE INDEX courses_tenant_slug_key ON courses(tenant_id, slug);

DROP INDEX IF EXISTS user_preferences_tenant_user_scope_key;
ALTER TABLE user_preferences ADD PRIMARY KEY (tenant_id, user_id, scope);
