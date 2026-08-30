package rubric

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Grade(c context.Context, s string, total int, feedback string, in []domain.RubricScore) (e error) {
	tx, e := r.db.BeginTx(c, nil)
	if e != nil {
		return
	}
	defer tx.Rollback()
	for _, x := range in {
		_, e = tx.ExecContext(c, `INSERT INTO submission_rubric_scores(submission_id,criterion_id,score,feedback) VALUES($1,$2,$3,$4) ON CONFLICT(submission_id,criterion_id) DO UPDATE SET score=EXCLUDED.score,feedback=EXCLUDED.feedback`, s, x.CriterionID, x.Score, x.Feedback)
		if e != nil {
			return
		}
	}
	_, e = tx.ExecContext(c, `UPDATE submissions SET score=$1,feedback=$2,status='graded',graded_at=NOW() WHERE id=$3`, total, feedback, s)
	if e != nil {
		return
	}
	return tx.Commit()
}
