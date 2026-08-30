package quiz

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, q domain.Quiz) (string, error) {
	query := `
		INSERT INTO quizzes (course_id, section_id, title, description, time_limit_secs, max_attempts, scheduled_at, available_until)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query,
		q.CourseID, q.SectionID, q.Title, q.Description, q.TimeLimitSecs, q.MaxAttempts, q.ScheduledAt, q.AvailableUntil,
	).Scan(&id)
	return id, err
}
