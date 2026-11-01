package rubric

import (
	"context"
	"errors"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

var ErrNotInTenant = errors.New("resource does not belong to this tenant")

var ErrSubmissionChanged = errors.New("submission changed while being graded")

var ErrAssignmentGraded = errors.New("assignment already has graded submissions")

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

	if _, e = tx.ExecContext(c, `UPDATE assignments SET max_score = $1 WHERE id = $2`, maxScore, a); e != nil {
		return
	}
	return tx.Commit()
}
