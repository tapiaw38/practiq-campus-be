package rubric

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) Grade(c context.Context, s string, total int, feedback string, in []domain.RubricScore, expectedVersion int) (e error) {
	guard, tenantID, e := tenantcontext.ExistsIn(c,
		"submissions sub JOIN assignments asg ON asg.id = sub.assignment_id JOIN courses co ON co.id = asg.course_id",
		"sub.id = $1", "co.tenant_id", 2)
	if e != nil {
		return
	}

	tx, e := r.db.BeginTx(c, nil)
	if e != nil {
		return
	}
	defer tx.Rollback()

	var belongs bool
	if e = tx.QueryRowContext(c, `SELECT `+guard, s, tenantID).Scan(&belongs); e != nil {
		return
	}
	if !belongs {
		return ErrNotInTenant
	}

	for _, x := range in {
		_, e = tx.ExecContext(c, `INSERT INTO submission_rubric_scores(submission_id,criterion_id,score,feedback) VALUES($1,$2,$3,$4) ON CONFLICT(submission_id,criterion_id) DO UPDATE SET score=EXCLUDED.score,feedback=EXCLUDED.feedback`, s, x.CriterionID, x.Score, x.Feedback)
		if e != nil {
			return
		}
	}
	res, e := tx.ExecContext(c, `UPDATE submissions SET score=$1,feedback=$2,status='graded',graded_at=NOW() WHERE id=$3 AND version=$4`, total, feedback, s, expectedVersion)
	if e != nil {
		return
	}
	if affected, err := res.RowsAffected(); err == nil && affected == 0 {
		return ErrSubmissionChanged
	}
	if e != nil {
		return
	}
	return tx.Commit()
}
