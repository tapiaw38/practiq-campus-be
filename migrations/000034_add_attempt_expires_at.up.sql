-- The deadline is pinned when the attempt starts, not derived from the quiz
-- at submit time: a teacher editing time_limit_secs mid-exam would otherwise
-- lengthen or cut short attempts already running.
-- NULL means untimed, which is what every attempt started before this column
-- existed was.
ALTER TABLE quiz_attempts ADD COLUMN expires_at TIMESTAMPTZ;
