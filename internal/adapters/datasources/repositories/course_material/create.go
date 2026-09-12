package course_material

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// ErrCourseNotInTenant is returned when the course is another institution's.
var ErrCourseNotInTenant = errors.New("course does not belong to this tenant")

func (r *repository) Create(ctx context.Context, m domain.CourseMaterial) (string, error) {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx, "courses c", "c.id = $1", "c.tenant_id", 9)
	if err != nil {
		return "", err
	}
	var id string
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO course_materials (course_id, assignment_id, section_id, uploader_id, title, description, kind, url)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8 WHERE `+guard+`
		RETURNING id`,
		m.CourseID, m.AssignmentID, m.SectionID, m.UploaderID, m.Title, m.Description, m.Kind, m.URL, tenantID,
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrCourseNotInTenant
	}
	return id, err
}
