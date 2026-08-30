package rubric

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) List(c context.Context, a string) (out []domain.RubricCriterion, e error) {
	out = make([]domain.RubricCriterion, 0)
	rows, e := r.db.QueryContext(c, `SELECT id,assignment_id,title,description,max_score,position FROM assignment_rubric_criteria WHERE assignment_id=$1 ORDER BY position`, a)
	if e != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var x domain.RubricCriterion
		e = rows.Scan(&x.ID, &x.AssignmentID, &x.Title, &x.Description, &x.MaxScore, &x.Position)
		if e != nil {
			return
		}
		out = append(out, x)
	}
	e = rows.Err()
	return
}
