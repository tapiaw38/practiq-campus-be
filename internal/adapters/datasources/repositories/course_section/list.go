package course_section

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.CourseSection, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "s.course_id").
		Where("s.course_id = ?", courseID).
		SQL(selectQualifiedSectionColumns, "course_sections s")
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query+" ORDER BY s.position ASC, s.created_at ASC", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []domain.CourseSection
	for rows.Next() {
		s, err := scanSection(rows)
		if err != nil {
			return nil, err
		}
		sections = append(sections, *s)
	}
	return sections, rows.Err()
}
