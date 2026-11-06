package submission

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) Grade(ctx context.Context, id string, score int, feedback string, expectedVersion int) error {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"assignments a JOIN courses c ON c.id = a.course_id",
		"a.id = submissions.assignment_id", "c.tenant_id", 6)
	if err != nil {
		return err
	}
	result, err := r.db.ExecContext(ctx, `
		UPDATE submissions SET score = $1, feedback = $2, status = $3, graded_at = NOW()
		WHERE id = $4 AND version = $5 AND `+guard,
		score, feedback, domain.SubmissionStatusGraded, id, expectedVersion, tenantID)
	if err != nil {
		return err
	}
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return ErrSubmissionChanged
	}
	return nil
}
