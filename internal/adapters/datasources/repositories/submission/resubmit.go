package submission

import (
	"context"
)

// Resubmit keeps one canonical delivery per student/activity while preserving
// teacher feedback only until the student explicitly sends a new version.
func (r *repository) Resubmit(ctx context.Context, id, content string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `DELETE FROM submission_rubric_scores WHERE submission_id = $1`, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE submissions SET content = $1, status = 'submitted', score = NULL, feedback = '', graded_at = NULL, submitted_at = NOW() WHERE id = $2`, content, id); err != nil {
		return err
	}
	return tx.Commit()
}
