package quiz

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// Qualified for the tenant join, which brings courses into scope. The question
// count stays a correlated subquery so a quiz with no questions still returns.
const selectQualifiedQuizColumns = `
	q.id, q.course_id, q.section_id, q.title, q.description, q.time_limit_secs, q.max_attempts,
	q.scheduled_at, q.available_until, q.created_at, q.updated_at,
	(SELECT COUNT(*) FROM quiz_questions WHERE quiz_id = q.id),
	q.weight, q.visible_group_id, q.unlock_after_type, q.unlock_after_id
`

func scanQuiz(row interface{ Scan(...any) error }) (*domain.Quiz, error) {
	var q domain.Quiz
	err := row.Scan(&q.ID, &q.CourseID, &q.SectionID, &q.Title, &q.Description, &q.TimeLimitSecs,
		&q.MaxAttempts, &q.ScheduledAt, &q.AvailableUntil, &q.CreatedAt, &q.UpdatedAt,
		&q.QuestionCount, &q.Weight, &q.VisibleGroupID, &q.UnlockAfterType, &q.UnlockAfterID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Quiz, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "q.course_id").
		Where("q.id = ?", id).
		SQL(selectQualifiedQuizColumns, "quizzes q")
	if err != nil {
		return nil, err
	}
	return scanQuiz(r.db.QueryRowContext(ctx, query, args...))
}
