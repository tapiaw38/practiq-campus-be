package course_section

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.CourseSection, error) {
	query := `SELECT id, course_id, title, position, created_at, updated_at FROM course_sections WHERE id = $1`
	var s domain.CourseSection
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.CourseID, &s.Title, &s.Position, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}
