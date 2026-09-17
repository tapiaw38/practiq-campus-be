package rubric

import (
	"context"
	"errors"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// ErrNotInTenant is returned when the assignment or submission belongs to
// another institution.
var ErrNotInTenant = errors.New("resource does not belong to this tenant")

// ErrSubmissionChanged means the student resubmitted while the teacher was
// grading, so the grade would have landed on content the teacher never read.
var ErrSubmissionChanged = errors.New("submission changed while being graded")

// ErrAssignmentGraded means grading has already begun, so neither the rubric
// nor the maximum it defines may move.
var ErrAssignmentGraded = errors.New("assignment already has graded submissions")

// Replace rewrites the rubric and, with it, the assignment's maximum score.
//
// The two are one fact — the criteria add up to the maximum — and they used to
// be written by two separate calls, so a failure between them left a rubric
// totalling one thing against a maximum saying another, with grading free to
// run in the gap. They go in the same transaction now, and the check that
// nothing has been graded yet is repeated inside it rather than trusted from
// a read taken beforehand.
//
// The tenant is checked before the delete: an assignment id from another
// institution would otherwise erase its criteria.
func (r *repository) Replace(c context.Context, a string, in []domain.RubricCriterion, maxScore int) (e error) {
	guard, tenantID, e := tenantcontext.ExistsIn(c,
		"assignments asg JOIN courses co ON co.id = asg.course_id", "asg.id = $1", "co.tenant_id", 2)
	if e != nil {
		return
	}

	tx, e := r.db.BeginTx(c, nil)
	if e != nil {
		return
	}
	defer tx.Rollback()

	var belongs bool
	if e = tx.QueryRowContext(c, `SELECT `+guard, a, tenantID).Scan(&belongs); e != nil {
		return
	}
	if !belongs {
		return ErrNotInTenant
	}

	// Locks this assignment's submissions for the transaction, so a grade
	// being written right now either lands first and blocks this, or waits
	// and then finds the rubric already replaced. Counted in Go because
	// Postgres does not allow FOR UPDATE alongside an aggregate.
	rows, e := tx.QueryContext(c, `SELECT graded_at FROM submissions WHERE assignment_id = $1 FOR UPDATE`, a)
	if e != nil {
		return
	}
	graded := 0
	for rows.Next() {
		var gradedAt *time.Time
		if e = rows.Scan(&gradedAt); e != nil {
			rows.Close()
			return
		}
		if gradedAt != nil {
			graded++
		}
	}
	if e = rows.Err(); e != nil {
		rows.Close()
		return
	}
	rows.Close()
	if graded > 0 {
		return ErrAssignmentGraded
	}

	if _, e = tx.ExecContext(c, `DELETE FROM assignment_rubric_criteria WHERE assignment_id=$1`, a); e != nil {
		return
	}
	for i, x := range in {
		if _, e = tx.ExecContext(c, `INSERT INTO assignment_rubric_criteria(assignment_id,title,description,max_score,position) VALUES($1,$2,$3,$4,$5)`, a, x.Title, x.Description, x.MaxScore, i); e != nil {
			return
		}
	}
	// The criteria define the maximum, so it moves with them — in the same
	// transaction, never as a second call that can fail on its own.
	if _, e = tx.ExecContext(c, `UPDATE assignments SET max_score = $1 WHERE id = $2`, maxScore, a); e != nil {
		return
	}
	return tx.Commit()
}
