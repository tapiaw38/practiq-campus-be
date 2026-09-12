package rubric

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) Scores(c context.Context, s string) (out []domain.RubricScore, e error) {
	tenantID := tenantcontext.ID(c)
	if tenantID == "" {
		return nil, tenantcontext.ErrNoTenant
	}
	rows, e := r.db.QueryContext(c, `
		SELECT sc.id,sc.submission_id,sc.criterion_id,sc.score,sc.feedback
		FROM submission_rubric_scores sc
		JOIN submissions sub ON sub.id = sc.submission_id
		JOIN assignments asg ON asg.id = sub.assignment_id
		JOIN courses co ON co.id = asg.course_id AND co.tenant_id = $2
		WHERE sc.submission_id=$1`, s, tenantID)
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
