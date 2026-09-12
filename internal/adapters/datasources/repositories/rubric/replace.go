package rubric

import (
	"context"
	"errors"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// ErrNotInTenant is returned when the assignment or submission belongs to
// another institution.
var ErrNotInTenant = errors.New("resource does not belong to this tenant")

// Replace clears the rubric before rewriting it, so the tenant is checked
// before the delete: an assignment id from another institution would otherwise
// erase its criteria.
func (r *repository) Replace(c context.Context, a string, in []domain.RubricCriterion) (e error) {
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

	if _, e = tx.ExecContext(c, `DELETE FROM assignment_rubric_criteria WHERE assignment_id=$1`, a); e != nil {
		return
	}
	for i, x := range in {
		if _, e = tx.ExecContext(c, `INSERT INTO assignment_rubric_criteria(assignment_id,title,description,max_score,position) VALUES($1,$2,$3,$4,$5)`, a, x.Title, x.Description, x.MaxScore, i); e != nil {
			return
		}
	}
	return tx.Commit()
}
