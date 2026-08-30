package enrollment

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

const selectEnrollmentColumns = `
	id, course_id, user_id, enrollment_role, status, enrolled_at
`

func scanEnrollment(row *sql.Row) (*domain.Enrollment, error) {
	var e domain.Enrollment
	err := row.Scan(&e.ID, &e.CourseID, &e.UserID, &e.EnrollmentRole, &e.Status, &e.EnrolledAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Enrollment, error) {
	row := r.db.QueryRowContext(ctx, "SELECT "+selectEnrollmentColumns+" FROM enrollments WHERE id = $1", id)
	return scanEnrollment(row)
}

func (r *repository) GetByCourseAndUser(ctx context.Context, courseID, userID string) (*domain.Enrollment, error) {
	row := r.db.QueryRowContext(
		ctx,
		"SELECT "+selectEnrollmentColumns+" FROM enrollments WHERE course_id = $1 AND user_id = $2",
		courseID, userID,
	)
	return scanEnrollment(row)
}
