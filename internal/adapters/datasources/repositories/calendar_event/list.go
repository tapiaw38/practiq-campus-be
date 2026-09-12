package calendar_event

import (
	"context"

	"github.com/lib/pq"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) ListByOwner(ctx context.Context, ownerID string) ([]domain.CalendarEvent, error) {
	tenantID := tenantcontext.ID(ctx)
	if tenantID == "" {
		return nil, tenantcontext.ErrNoTenant
	}
	query := `
		SELECT ` + selectEventColumns + `
		FROM calendar_events
		WHERE tenant_id = $2 AND (owner_id = $1 OR EXISTS (
			SELECT 1 FROM calendar_event_attendees attendee
			WHERE attendee.event_id = calendar_events.id AND attendee.user_id = $1
		))
		ORDER BY starts_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, ownerID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.CalendarEvent
	for rows.Next() {
		var e domain.CalendarEvent
		if err := rows.Scan(&e.ID, &e.OwnerID, &e.CourseID, &e.Title, &e.Description, &e.StartsAt, &e.EndsAt, &e.AllDay, &e.RecurrenceRule, &e.ReminderMinutes, &e.CreatedAt, pq.Array(&e.AttendeeIDs)); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}
