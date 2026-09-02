package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
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
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT "+selectEnrollmentColumns+" FROM enrollments WHERE course_id = $1 ORDER BY enrolled_at DESC",
		courseID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEnrollmentRows(rows)
}

func (r *repository) ListByUser(ctx context.Context, userID string) ([]domain.Enrollment, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT "+selectEnrollmentColumns+" FROM enrollments WHERE user_id = $1 ORDER BY enrolled_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEnrollmentRows(rows)
}
