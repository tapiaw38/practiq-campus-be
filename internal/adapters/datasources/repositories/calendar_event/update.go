package calendar_event

import (
	"context"
	"fmt"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, event domain.CalendarEvent) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `UPDATE calendar_events SET course_id=$2, title=$3, description=$4, starts_at=$5, ends_at=$6, all_day=$7, recurrence_rule=$8, reminder_minutes=$9 WHERE id=$1`, id, event.CourseID, event.Title, event.Description, event.StartsAt, event.EndsAt, event.AllDay, event.RecurrenceRule, event.ReminderMinutes)
	if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM calendar_event_attendees WHERE event_id=$1`, id); err != nil {
		return err
	}
	for _, userID := range event.AttendeeIDs {
		if _, err = tx.ExecContext(ctx, `INSERT INTO calendar_event_attendees (event_id, user_id) VALUES ($1, $2)`, id, userID); err != nil {
			return fmt.Errorf("update attendee: %w", err)
		}
	}
	return tx.Commit()
}
