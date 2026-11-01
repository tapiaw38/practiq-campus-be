package domain

import "time"

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
