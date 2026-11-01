package submission

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

var ErrSubmissionGraded = errors.New("submission has already been graded")

var ErrSubmissionGone = errors.New("submission no longer exists")

var ErrSubmissionChanged = errors.New("submission changed while being graded")

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

	var gradedAt *time.Time
	err = tx.QueryRowContext(ctx,
		`SELECT graded_at FROM submissions WHERE id = $1 FOR UPDATE`, id).Scan(&gradedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrSubmissionGone
	}
	if err != nil {
		return err
	}
	if gradedAt != nil {
		return ErrSubmissionGraded
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM submission_rubric_scores WHERE submission_id = $1`, id); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE submissions
		SET content = $1, status = 'submitted', score = NULL, feedback = '', graded_at = NULL,
		    submitted_at = NOW(), version = version + 1
		WHERE id = $2 AND graded_at IS NULL`, content, id)
	if err != nil {
		return err
	}
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return ErrSubmissionGraded
	}
	return tx.Commit()
}
