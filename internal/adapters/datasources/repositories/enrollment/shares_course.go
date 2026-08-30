package enrollment

import "context"

func (r *repository) SharesCourseWith(ctx context.Context, userA, userB string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM courses c
			JOIN enrollments e ON e.course_id = c.id AND e.status = 'active'
			WHERE (c.owner_id = $1 AND e.user_id = $2) OR (c.owner_id = $2 AND e.user_id = $1)
		) OR EXISTS (
			SELECT 1 FROM enrollments e1
			JOIN enrollments e2 ON e2.course_id = e1.course_id
			WHERE e1.user_id = $1 AND e2.user_id = $2 AND e1.status = 'active' AND e2.status = 'active'
		)
	`
	var shares bool
	err := r.db.QueryRowContext(ctx, query, userA, userB).Scan(&shares)
	return shares, err
}
