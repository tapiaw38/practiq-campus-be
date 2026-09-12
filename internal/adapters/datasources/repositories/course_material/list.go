package course_material

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.CourseMaterial, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "m.course_id").
		Where("m.course_id = ?", courseID).
		SQL(selectQualifiedMaterialColumns, "course_materials m")
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query+" ORDER BY m.created_at DESC", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []domain.CourseMaterial
	for rows.Next() {
		m, err := scanMaterial(rows)
		if err != nil {
			return nil, err
		}
		materials = append(materials, *m)
	}
	return materials, rows.Err()
}
