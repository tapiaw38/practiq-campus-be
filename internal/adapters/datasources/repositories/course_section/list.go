package course_section

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.CourseSection, error) {
	query := `
		SELECT id, course_id, title, position, created_at, updated_at
		FROM course_sections
		WHERE course_id = $1
		ORDER BY position ASC, created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []domain.CourseSection
	for rows.Next() {
		var s domain.CourseSection
		if err := rows.Scan(&s.ID, &s.CourseID, &s.Title, &s.Position, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sections = append(sections, s)
	}
	return sections, rows.Err()
}
