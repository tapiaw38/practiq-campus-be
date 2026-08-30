package course

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func scanCourse(row *sql.Row) (*domain.Course, error) {
	var c domain.Course
	err := row.Scan(
		&c.ID, &c.OwnerID, &c.Title, &c.Slug, &c.Description, &c.Status,
		&c.StartDate, &c.EndDate, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

const selectCourseColumns = `
	id, owner_id, title, slug, description, status, start_date, end_date, created_at, updated_at
`

func (r *repository) Get(ctx context.Context, id string) (*domain.Course, error) {
	row := r.db.QueryRowContext(ctx, "SELECT "+selectCourseColumns+" FROM courses WHERE id = $1", id)
	return scanCourse(row)
}

func (r *repository) GetBySlug(ctx context.Context, slug string) (*domain.Course, error) {
	row := r.db.QueryRowContext(ctx, "SELECT "+selectCourseColumns+" FROM courses WHERE slug = $1", slug)
	return scanCourse(row)
}
