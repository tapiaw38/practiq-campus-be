package rubric

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Replace(c context.Context, a string, in []domain.RubricCriterion) (e error) {
	tx, e := r.db.BeginTx(c, nil)
	if e != nil {
		return
	}
	defer tx.Rollback()
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
