package rubric

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Scores(c context.Context, s string) (out []domain.RubricScore, e error) {
	rows, e := r.db.QueryContext(c, `SELECT id,submission_id,criterion_id,score,feedback FROM submission_rubric_scores WHERE submission_id=$1`, s)
	if e != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var x domain.RubricScore
		e = rows.Scan(&x.ID, &x.SubmissionID, &x.CriterionID, &x.Score, &x.Feedback)
		if e != nil {
			return
		}
		out = append(out, x)
	}
	e = rows.Err()
	return
}
