package quiz

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, q domain.Quiz) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE quizzes SET section_id=$1, title=$2, description=$3, time_limit_secs=$4, max_attempts=$5, scheduled_at=$6, available_until=$7, updated_at=now()
		WHERE id=$8
	`, q.SectionID, q.Title, q.Description, q.TimeLimitSecs, q.MaxAttempts, q.ScheduledAt, q.AvailableUntil, id)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM quizzes WHERE id=$1`, id)
	return err
}
