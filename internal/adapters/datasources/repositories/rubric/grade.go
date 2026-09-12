package rubric

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// Grade writes marks onto a student's work, so the submission is confirmed to
// belong to this institution before anything is written.
func (r *repository) Grade(c context.Context, s string, total int, feedback string, in []domain.RubricScore) (e error) {
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
	_, e = tx.ExecContext(c, `UPDATE submissions SET score=$1,feedback=$2,status='graded',graded_at=NOW() WHERE id=$3`, total, feedback, s)
	if e != nil {
		return
	}
	return tx.Commit()
}
