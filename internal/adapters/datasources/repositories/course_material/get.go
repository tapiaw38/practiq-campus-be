package course_material

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

const selectQualifiedMaterialColumns = `
	m.id, m.course_id, m.assignment_id, m.section_id, m.uploader_id, m.title, m.description, m.kind, m.url, m.created_at
`

func scanMaterial(row interface{ Scan(...any) error }) (*domain.CourseMaterial, error) {
	var m domain.CourseMaterial
	err := row.Scan(&m.ID, &m.CourseID, &m.AssignmentID, &m.SectionID, &m.UploaderID,
		&m.Title, &m.Description, &m.Kind, &m.URL, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.CourseMaterial, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "m.course_id").
		Where("m.id = ?", id).
		SQL(selectQualifiedMaterialColumns, "course_materials m")
	if err != nil {
		return nil, err
	}
	return scanMaterial(r.db.QueryRowContext(ctx, query, args...))
}
