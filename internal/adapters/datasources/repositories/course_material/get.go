package course_material

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

const selectMaterialColumns = `
	id, course_id, assignment_id, section_id, uploader_id, title, description, kind, url, created_at
`

func (r *repository) Get(ctx context.Context, id string) (*domain.CourseMaterial, error) {
	row := r.db.QueryRowContext(ctx, "SELECT "+selectMaterialColumns+" FROM course_materials WHERE id = $1", id)

	var m domain.CourseMaterial
	err := row.Scan(&m.ID, &m.CourseID, &m.AssignmentID, &m.SectionID, &m.UploaderID, &m.Title, &m.Description, &m.Kind, &m.URL, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}
