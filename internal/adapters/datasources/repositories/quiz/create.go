package quiz

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// ErrCourseNotInTenant is returned when the course is another institution's.
var ErrCourseNotInTenant = errors.New("course does not belong to this tenant")

func (r *repository) Create(ctx context.Context, q domain.Quiz) (string, error) {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx, "courses c", "c.id = $1", "c.tenant_id", 13)
	if err != nil {
		return "", err
	}
	var id string
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO quizzes (course_id, section_id, title, description, time_limit_secs, max_attempts, scheduled_at, available_until, weight, visible_group_id, unlock_after_type, unlock_after_id)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12 WHERE `+guard+`
		RETURNING id`,
		q.CourseID, q.SectionID, q.Title, q.Description, q.TimeLimitSecs, q.MaxAttempts,
		q.ScheduledAt, q.AvailableUntil, q.Weight, q.VisibleGroupID, q.UnlockAfterType, q.UnlockAfterID, tenantID,
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrCourseNotInTenant
	}
	return id, err
}
