package course_material

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, m domain.CourseMaterial) (string, error) {
	query := `
		INSERT INTO course_materials (course_id, section_id, uploader_id, title, description, kind, url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query,
		m.CourseID, m.SectionID, m.UploaderID, m.Title, m.Description, m.Kind, m.URL,
	).Scan(&id)
	return id, err
}
