package quiz

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, q domain.Quiz) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE quizzes SET section_id=$1, title=$2, description=$3, time_limit_secs=$4, max_attempts=$5, scheduled_at=$6, available_until=$7, weight=$8, visible_group_id=$9, unlock_after_type=$10, unlock_after_id=$11, updated_at=now()
		WHERE id=$12
	`, q.SectionID, q.Title, q.Description, q.TimeLimitSecs, q.MaxAttempts, q.ScheduledAt, q.AvailableUntil, q.Weight, q.VisibleGroupID, q.UnlockAfterType, q.UnlockAfterID, id)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM quizzes WHERE id=$1`, id)
	return err
}
