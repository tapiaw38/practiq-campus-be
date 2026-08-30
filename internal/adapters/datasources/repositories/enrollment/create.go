package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, e domain.Enrollment) (string, error) {
	query := `
		INSERT INTO enrollments (course_id, user_id, enrollment_role, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, e.CourseID, e.UserID, e.EnrollmentRole, e.Status).Scan(&id)
	return id, err
}
