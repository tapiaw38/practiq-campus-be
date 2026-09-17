package submission

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// ErrSubmissionGraded means the work was graded before this resubmission got
// to it. Resubmitting clears the grade, so it must not proceed on a stale
// read of that state.
var ErrSubmissionGraded = errors.New("submission has already been graded")

// ErrSubmissionGone means the submission disappeared between being read and
// being locked.
var ErrSubmissionGone = errors.New("submission no longer exists")

// ErrSubmissionChanged means the student resubmitted while the teacher was
// grading, so the grade would have landed on content the teacher never read.
var ErrSubmissionChanged = errors.New("submission changed while being graded")

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

	// Whether the work has been graded is re-read here, under a row lock,
	// because this wipes the grading. The caller checks it too, but that read
	// and this delete are separate moments: a teacher committing a grade in
	// between would have it erased by a resubmission that was allowed on the
	// strength of a stale "not graded yet".
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
	// The same condition again in the write: if a grade landed between the
	// lock being taken and this statement — impossible while the lock holds,
	// but free to assert — the update simply does not match.
	// version moves with the content, so a grade written against the version
	// the teacher read will no longer match once this lands.
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
