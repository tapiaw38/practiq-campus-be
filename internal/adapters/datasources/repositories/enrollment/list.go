package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func scanEnrollmentRows(rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}) ([]domain.Enrollment, error) {
	var enrollments []domain.Enrollment
	for rows.Next() {
		var e domain.Enrollment
		if err := rows.Scan(&e.ID, &e.CourseID, &e.UserID, &e.EnrollmentRole, &e.Status, &e.EnrolledAt); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, e)
	}
	return enrollments, rows.Err()
}

func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.Enrollment, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "e.course_id").
		Where("e.course_id = ?", courseID).
		SQL(selectQualifiedEnrollmentColumns, "enrollments e")
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query+" ORDER BY e.enrolled_at DESC", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEnrollmentRows(rows)
}

// ListByUser answers for one institution. A student enrolled in two sees the
// courses of the one they are looking at, not both lists merged.
func (r *repository) ListByUser(ctx context.Context, userID string) ([]domain.Enrollment, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "e.course_id").
		Where("e.user_id = ?", userID).
		SQL(selectQualifiedEnrollmentColumns, "enrollments e")
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query+" ORDER BY e.enrolled_at DESC", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEnrollmentRows(rows)
}
