package dashboard

import (
	"context"
	"database/sql"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

type CourseSummary struct {
	CourseID       string
	CompletedItems int
	TotalItems     int
	PendingItems   int
	NextDueAt      *time.Time
	NextDueTitle   string
}

type Repository interface {
	StudentCourses(context.Context, string) ([]CourseSummary, error)
}
type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

// StudentCourses aggregates assignments and quizzes in one tenant-scoped query.
// Keeping this server-side prevents the dashboard from issuing one request per
// course and per activity merely to draw its course cards.
func (r *repository) StudentCourses(ctx context.Context, userID string) ([]CourseSummary, error) {
	tenantID := tenantcontext.ID(ctx)
	if tenantID == "" {
		return nil, tenantcontext.ErrNoTenant
	}
	rows, err := r.db.QueryContext(ctx, `
		WITH mine AS (
			SELECT c.id FROM courses c JOIN enrollments e ON e.course_id=c.id
			WHERE c.tenant_id=$1 AND e.user_id=$2
		), items AS (
			SELECT a.course_id, a.id, a.title, a.due_at, EXISTS(
				SELECT 1 FROM submissions s WHERE s.assignment_id=a.id AND s.user_id=$2
			) AS completed FROM assignments a JOIN mine ON mine.id=a.course_id
			UNION ALL
			SELECT q.course_id, q.id, q.title, q.available_until, EXISTS(
				SELECT 1 FROM quiz_attempts qa WHERE qa.quiz_id=q.id AND qa.user_id=$2 AND qa.submitted_at IS NOT NULL
			) AS completed FROM quizzes q JOIN mine ON mine.id=q.course_id
		)
		SELECT mine.id, COUNT(items.id), COUNT(items.id) FILTER (WHERE items.completed),
			COUNT(items.id) FILTER (WHERE NOT items.completed),
			MIN(items.due_at) FILTER (WHERE NOT items.completed AND items.due_at >= NOW()),
			COALESCE((ARRAY_AGG(items.title ORDER BY items.due_at) FILTER (WHERE NOT items.completed AND items.due_at >= NOW()))[1], '')
		FROM mine LEFT JOIN items ON items.course_id=mine.id GROUP BY mine.id`, tenantID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CourseSummary
	for rows.Next() {
		var item CourseSummary
		if err := rows.Scan(&item.CourseID, &item.TotalItems, &item.CompletedItems, &item.PendingItems, &item.NextDueAt, &item.NextDueTitle); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}
