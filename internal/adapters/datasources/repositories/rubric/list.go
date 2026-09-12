package rubric

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// Criteria hang off an assignment, which belongs to a tenant through its
// course.
func (r *repository) List(c context.Context, a string) (out []domain.RubricCriterion, e error) {
	out = make([]domain.RubricCriterion, 0)
	tenantID := tenantcontext.ID(c)
	if tenantID == "" {
		return nil, tenantcontext.ErrNoTenant
	}
	rows, e := r.db.QueryContext(c, `
		SELECT rc.id,rc.assignment_id,rc.title,rc.description,rc.max_score,rc.position
		FROM assignment_rubric_criteria rc
		JOIN assignments a ON a.id = rc.assignment_id
		JOIN courses co ON co.id = a.course_id AND co.tenant_id = $2
		WHERE rc.assignment_id=$1 ORDER BY rc.position`, a, tenantID)
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
