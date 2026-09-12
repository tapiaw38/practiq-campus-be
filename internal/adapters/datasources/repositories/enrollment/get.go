package enrollment

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// Qualified for the tenant join, which brings courses into scope.
const selectQualifiedEnrollmentColumns = `
	e.id, e.course_id, e.user_id, e.enrollment_role, e.status, e.enrolled_at
`

func scanEnrollment(row interface{ Scan(...any) error }) (*domain.Enrollment, error) {
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
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "e.course_id").
		Where("e.id = ?", id).
		SQL(selectQualifiedEnrollmentColumns, "enrollments e")
	if err != nil {
		return nil, err
	}
	return scanEnrollment(r.db.QueryRowContext(ctx, query, args...))
}

func (r *repository) GetByCourseAndUser(ctx context.Context, courseID, userID string) (*domain.Enrollment, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "e.course_id").
		Where("e.course_id = ?", courseID).
		Where("e.user_id = ?", userID).
		SQL(selectQualifiedEnrollmentColumns, "enrollments e")
	if err != nil {
		return nil, err
	}
	return scanEnrollment(r.db.QueryRowContext(ctx, query, args...))
}
