package course_section

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// Qualified for the tenant join, which brings courses into scope.
const selectQualifiedSectionColumns = `
	s.id, s.course_id, s.title, s.description, s.position, s.created_at, s.updated_at
`

func scanSection(row interface{ Scan(...any) error }) (*domain.CourseSection, error) {
	var s domain.CourseSection
	err := row.Scan(&s.ID, &s.CourseID, &s.Title, &s.Description, &s.Position, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Sections belong to a tenant through their course, so an id from another
// institution does not resolve.
func (r *repository) Get(ctx context.Context, id string) (*domain.CourseSection, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "s.course_id").
		Where("s.id = ?", id).
		SQL(selectQualifiedSectionColumns, "course_sections s")
	if err != nil {
		return nil, err
	}
	return scanSection(r.db.QueryRowContext(ctx, query, args...))
}
