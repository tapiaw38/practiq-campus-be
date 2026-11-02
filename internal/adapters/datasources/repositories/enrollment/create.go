package enrollment

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

var ErrCourseNotInTenant = errors.New("course does not belong to this tenant")

func (r *repository) Create(ctx context.Context, e domain.Enrollment) (string, error) {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx, "courses c", "c.id = $1", "c.tenant_id", 5)
	if err != nil {
		return "", err
	}
	var id string
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO enrollments (course_id, user_id, enrollment_role, status)
		SELECT $1, $2, $3, $4 WHERE `+guard+`
		RETURNING id`, e.CourseID, e.UserID, e.EnrollmentRole, e.Status, tenantID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrCourseNotInTenant
	}
	return id, err
}
