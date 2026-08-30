package domain

import "time"

// CalendarEvent is a manually-created personal event only. Assignment due
// dates are never copied in here — they're computed at read time from
// Assignment.DueAt, so there's exactly one place that can go stale.
type CalendarEvent struct {
	ID              string
	OwnerID         string
	CourseID        *string
	Title           string
	Description     string
	StartsAt        time.Time
	EndsAt          *time.Time
	AllDay          bool
	RecurrenceRule  string
	ReminderMinutes *int
	AttendeeIDs     []string
	CreatedAt       time.Time
}
