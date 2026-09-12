-- Reverting this is only possible while no one keeps preferences in more than
-- one institution. The old key was (user_id, scope); once the same person has
-- a row per institution, restoring it means choosing which of their settings
-- to destroy, and a migration is the wrong place to make that choice
-- silently.
--
-- So it is checked first and refused with a readable message, rather than
-- failing on a unique_violation that says nothing about what to do.
DO $$
DECLARE
  duplicates INT;
BEGIN
  SELECT count(*) INTO duplicates FROM (
    SELECT user_id, scope FROM user_preferences
    GROUP BY user_id, scope HAVING count(*) > 1
  ) AS d;

  IF duplicates > 0 THEN
    RAISE EXCEPTION
      'cannot revert: % user/scope pairs hold preferences in more than one institution. '
      'Restoring the old primary key would discard all but one of each. '
      'Decide which institution''s preferences to keep and delete the rest first.',
      duplicates;
  END IF;
END $$;

DROP INDEX IF EXISTS idx_user_preferences_tenant_user;
DROP INDEX IF EXISTS user_preferences_tenant_user_scope_key;
ALTER TABLE user_preferences DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE user_preferences ADD PRIMARY KEY (user_id, scope);
