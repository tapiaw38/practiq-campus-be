CREATE TABLE IF NOT EXISTS calendar_event_attendees (
  event_id UUID NOT NULL REFERENCES calendar_events(id) ON DELETE CASCADE,
  user_id VARCHAR(255) NOT NULL REFERENCES campus_profiles(id) ON DELETE CASCADE,
  PRIMARY KEY (event_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_calendar_event_attendees_user_id ON calendar_event_attendees(user_id);
