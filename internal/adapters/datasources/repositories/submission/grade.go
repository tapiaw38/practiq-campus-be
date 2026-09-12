package submission

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// Grade writes a mark onto a student's work, so it is guarded by the same
// chain the reads use: a submission id from another institution updates
// nothing.
func (r *repository) Grade(ctx context.Context, id string, score int, feedback string) error {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"assignments a JOIN courses c ON c.id = a.course_id",
		"a.id = submissions.assignment_id", "c.tenant_id", 5)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		UPDATE submissions SET score = $1, feedback = $2, status = $3, graded_at = NOW()
		WHERE id = $4 AND `+guard,
		score, feedback, domain.SubmissionStatusGraded, id, tenantID)
	return err
}
