package calendar_event

import (
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type EventData struct {
	ID              string   `json:"id"`
	OwnerID         string   `json:"owner_id,omitempty"`
	CourseID        *string  `json:"course_id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	StartsAt        string   `json:"starts_at"`
	EndsAt          *string  `json:"ends_at"`
	AllDay          bool     `json:"all_day"`
	RecurrenceRule  string   `json:"recurrence_rule"`
	ReminderMinutes *int     `json:"reminder_minutes"`
	AttendeeIDs     []string `json:"attendee_ids,omitempty"`
	// Source is "manual" for a user-created event or "assignment_due" for a
	// synthetic entry computed from an assignment's due date — the latter
	// is never persisted, only assembled at read time.
	Source string `json:"source"`
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02T15:04:05Z")
}

func formatOptionalTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := formatTime(*t)
	return &s
}

func toEventData(e domain.CalendarEvent) EventData {
	return EventData{
		ID:              e.ID,
		OwnerID:         e.OwnerID,
		CourseID:        e.CourseID,
		Title:           e.Title,
		Description:     e.Description,
		StartsAt:        formatTime(e.StartsAt),
		EndsAt:          formatOptionalTime(e.EndsAt),
		AllDay:          e.AllDay,
		RecurrenceRule:  e.RecurrenceRule,
		ReminderMinutes: e.ReminderMinutes,
		AttendeeIDs:     e.AttendeeIDs,
		Source:          "manual",
	}
}
