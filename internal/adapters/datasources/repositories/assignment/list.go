package assignment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// ListByCourse is scoped by the course's tenant as well as by the course, so a
// course id belonging to another institution returns nothing rather than its
// assignments.
func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.Assignment, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "a.course_id").
		Where("a.course_id = ?", courseID).
		SQL(selectQualifiedAssignmentColumns, "assignments a")
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query+" ORDER BY a.due_at ASC NULLS LAST, a.created_at ASC", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []domain.Assignment
	for rows.Next() {
		a, err := scanAssignment(rows)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, *a)
	}
	return assignments, rows.Err()
}
