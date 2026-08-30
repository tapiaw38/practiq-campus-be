package course_material

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.CourseMaterial, error) {
	query := `SELECT ` + selectMaterialColumns + ` FROM course_materials WHERE course_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []domain.CourseMaterial
	for rows.Next() {
		var m domain.CourseMaterial
		if err := rows.Scan(&m.ID, &m.CourseID, &m.AssignmentID, &m.SectionID, &m.UploaderID, &m.Title, &m.Description, &m.Kind, &m.URL, &m.CreatedAt); err != nil {
			return nil, err
		}
		materials = append(materials, m)
	}
	return materials, rows.Err()
}
