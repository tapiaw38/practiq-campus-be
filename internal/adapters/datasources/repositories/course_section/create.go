package course_section

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, s domain.CourseSection) (string, error) {
	query := `
		INSERT INTO course_sections (course_id, title, position)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, s.CourseID, s.Title, s.Position).Scan(&id)
	return id, err
}
