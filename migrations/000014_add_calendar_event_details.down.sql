ALTER TABLE calendar_events
  DROP CONSTRAINT IF EXISTS calendar_events_reminder_minutes_check,
  DROP CONSTRAINT IF EXISTS calendar_events_recurrence_rule_check,
  DROP COLUMN IF EXISTS reminder_minutes,
  DROP COLUMN IF EXISTS recurrence_rule,
  DROP COLUMN IF EXISTS description;
