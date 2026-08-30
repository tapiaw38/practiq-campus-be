ALTER TABLE calendar_events
  ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS recurrence_rule VARCHAR(16) NOT NULL DEFAULT 'none',
  ADD COLUMN IF NOT EXISTS reminder_minutes INTEGER;

ALTER TABLE calendar_events
  ADD CONSTRAINT calendar_events_recurrence_rule_check
  CHECK (recurrence_rule IN ('none', 'daily', 'weekly', 'monthly'));

ALTER TABLE calendar_events
  ADD CONSTRAINT calendar_events_reminder_minutes_check
  CHECK (reminder_minutes IS NULL OR reminder_minutes BETWEEN 0 AND 10080);
