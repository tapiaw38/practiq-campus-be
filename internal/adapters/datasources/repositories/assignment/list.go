package assignment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.Assignment, error) {
	query := `SELECT ` + selectAssignmentColumns + ` FROM assignments WHERE course_id = $1 ORDER BY due_at ASC NULLS LAST, created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []domain.Assignment
	for rows.Next() {
		var a domain.Assignment
		if err := rows.Scan(&a.ID, &a.CourseID, &a.SectionID, &a.Title, &a.Description, &a.DueAt, &a.MaxScore, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, rows.Err()
}
