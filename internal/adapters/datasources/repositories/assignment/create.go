package assignment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, a domain.Assignment) (string, error) {
	query := `
		INSERT INTO assignments (course_id, section_id, title, description, due_at, max_score)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query,
		a.CourseID, a.SectionID, a.Title, a.Description, a.DueAt, a.MaxScore,
	).Scan(&id)
	return id, err
}
