package calendar_event

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

const selectEventColumns = `id, owner_id, course_id, title, description, starts_at, ends_at, all_day, recurrence_rule, reminder_minutes, created_at,
  COALESCE((SELECT array_agg(user_id) FROM calendar_event_attendees WHERE event_id = calendar_events.id), '{}')`

func (r *repository) Get(ctx context.Context, id string) (*domain.CalendarEvent, error) {
	var event domain.CalendarEvent
	err := r.db.QueryRowContext(ctx, `SELECT `+selectEventColumns+` FROM calendar_events WHERE id = $1`, id).Scan(
		&event.ID, &event.OwnerID, &event.CourseID, &event.Title, &event.Description, &event.StartsAt, &event.EndsAt,
		&event.AllDay, &event.RecurrenceRule, &event.ReminderMinutes, &event.CreatedAt, pq.Array(&event.AttendeeIDs),
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &event, nil
}
