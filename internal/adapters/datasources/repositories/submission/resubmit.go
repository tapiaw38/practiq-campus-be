package submission

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// Resubmit keeps one canonical delivery per student/activity while preserving
// teacher feedback only until the student explicitly sends a new version.
//
// It clears rubric scores first, so the tenant has to be checked before that
// delete: an id from another institution would otherwise erase its grading.
func (r *repository) Resubmit(ctx context.Context, id, content string) error {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"assignments a JOIN courses c ON c.id = a.course_id",
		"a.id = (SELECT assignment_id FROM submissions WHERE id = $1)", "c.tenant_id", 2)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var belongs bool
	if err = tx.QueryRowContext(ctx, `SELECT `+guard, id, tenantID).Scan(&belongs); err != nil {
		return err
	}
	if !belongs {
		return ErrAssignmentNotInTenant
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM submission_rubric_scores WHERE submission_id = $1`, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE submissions SET content = $1, status = 'submitted', score = NULL, feedback = '', graded_at = NULL, submitted_at = NOW() WHERE id = $2`, content, id); err != nil {
		return err
	}
	return tx.Commit()
}
