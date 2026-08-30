package calendar_event

import (
	"context"
	"fmt"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, e domain.CalendarEvent) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	query := `
		INSERT INTO calendar_events (owner_id, course_id, title, description, starts_at, ends_at, all_day, recurrence_rule, reminder_minutes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	var id string
	err = tx.QueryRowContext(ctx, query, e.OwnerID, e.CourseID, e.Title, e.Description, e.StartsAt, e.EndsAt, e.AllDay, e.RecurrenceRule, e.ReminderMinutes).Scan(&id)
	if err != nil {
		return "", err
	}
	for _, userID := range e.AttendeeIDs {
		if _, err = tx.ExecContext(ctx, `INSERT INTO calendar_event_attendees (event_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, id, userID); err != nil {
			return "", fmt.Errorf("create attendee: %w", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return "", err
	}
	return id, nil
}
