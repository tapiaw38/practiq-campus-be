-- Reverting is only possible while no two institutions share a course slug.
-- The old constraint was UNIQUE(slug) across the whole platform; once two
-- institutions each have a "matematica", restoring it means deleting one of
-- their courses, which a migration must not decide on its own.
--
-- Checked before anything is dropped, so a refusal leaves the schema exactly
-- as it was. Doing the drops first and discovering the conflict at the end
-- left the tables without their columns and without their keys.
DO $$
DECLARE
  clashes INT;
BEGIN
  SELECT count(*) INTO clashes FROM (
    SELECT slug FROM courses GROUP BY slug HAVING count(*) > 1
  ) AS d;

  IF clashes > 0 THEN
    RAISE EXCEPTION
      'cannot revert: % course slugs are used by more than one institution. '
      'A platform-wide unique slug would require deleting one course of each pair. '
      'Rename or remove the duplicates first.',
      clashes;
  END IF;
END $$;

DROP INDEX IF EXISTS idx_notifications_tenant_user_created;
ALTER TABLE notifications DROP COLUMN IF EXISTS tenant_id;
DROP INDEX IF EXISTS idx_conversations_tenant_id;
ALTER TABLE conversations DROP COLUMN IF EXISTS tenant_id;
DROP INDEX IF EXISTS idx_calendar_events_tenant_id;
ALTER TABLE calendar_events DROP COLUMN IF EXISTS tenant_id;
DROP INDEX IF EXISTS courses_tenant_slug_key;
ALTER TABLE courses ADD CONSTRAINT courses_slug_key UNIQUE (slug);
DROP INDEX IF EXISTS idx_courses_tenant_id;
ALTER TABLE courses DROP COLUMN IF EXISTS tenant_id;
DROP TABLE IF EXISTS campus_tenants;
