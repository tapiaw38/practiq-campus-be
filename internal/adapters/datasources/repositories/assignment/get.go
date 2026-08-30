package assignment

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

const selectAssignmentColumns = `
	id, course_id, section_id, title, description, due_at, max_score, created_at, updated_at
`

func (r *repository) Get(ctx context.Context, id string) (*domain.Assignment, error) {
	row := r.db.QueryRowContext(ctx, "SELECT "+selectAssignmentColumns+" FROM assignments WHERE id = $1", id)

	var a domain.Assignment
	err := row.Scan(&a.ID, &a.CourseID, &a.SectionID, &a.Title, &a.Description, &a.DueAt, &a.MaxScore, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}
