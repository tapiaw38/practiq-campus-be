package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// SharesCourseWith decides who may open a conversation with whom, so it
// answers within one institution.
//
// Sharing a course in one institution is not standing to write to somebody in
// another. Unscoped, two people who teach together in one of them could open a
// thread from the other, where they have no relationship at all.
func (r *repository) SharesCourseWith(ctx context.Context, userA, userB string) (bool, error) {
	tenantID := tenantcontext.ID(ctx)
	if tenantID == "" {
		return false, tenantcontext.ErrNoTenant
	}
	query := `
		SELECT EXISTS (
			SELECT 1 FROM courses c
			JOIN enrollments e ON e.course_id = c.id AND e.status = 'active'
			WHERE c.tenant_id = $3
			AND ((c.owner_id = $1 AND e.user_id = $2) OR (c.owner_id = $2 AND e.user_id = $1))
		) OR EXISTS (
			SELECT 1 FROM enrollments e1
			JOIN enrollments e2 ON e2.course_id = e1.course_id
			JOIN courses c ON c.id = e1.course_id AND c.tenant_id = $3
			WHERE e1.user_id = $1 AND e2.user_id = $2 AND e1.status = 'active' AND e2.status = 'active'
		)
	`
	var shares bool
	err := r.db.QueryRowContext(ctx, query, userA, userB, tenantID).Scan(&shares)
	return shares, err
}
