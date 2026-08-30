package submission

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Grade(ctx context.Context, id string, score int, feedback string) error {
	query := `
		UPDATE submissions SET
			score = $1, feedback = $2, status = $3, graded_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, score, feedback, domain.SubmissionStatusGraded, id)
	return err
}
