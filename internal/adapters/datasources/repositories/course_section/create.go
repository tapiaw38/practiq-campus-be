package course_section

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// ErrCourseNotInTenant is returned when the course is another institution's.
var ErrCourseNotInTenant = errors.New("course does not belong to this tenant")

func (r *repository) Create(ctx context.Context, s domain.CourseSection) (string, error) {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx, "courses c", "c.id = $1", "c.tenant_id", 5)
	if err != nil {
		return "", err
	}
	var id string
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO course_sections (course_id, title, description, position)
		SELECT $1, $2, $3, $4 WHERE `+guard+`
		RETURNING id`, s.CourseID, s.Title, s.Description, s.Position, tenantID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrCourseNotInTenant
	}
	return id, err
}
