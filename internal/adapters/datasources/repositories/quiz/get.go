package quiz

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

const selectQuizColumns = `
	id, course_id, section_id, title, description, time_limit_secs, max_attempts, scheduled_at, available_until, created_at, updated_at,
	(SELECT COUNT(*) FROM quiz_questions WHERE quiz_id = quizzes.id),
	weight, visible_group_id, unlock_after_type, unlock_after_id
`

func scanQuiz(row *sql.Row) (*domain.Quiz, error) {
	var q domain.Quiz
	err := row.Scan(&q.ID, &q.CourseID, &q.SectionID, &q.Title, &q.Description, &q.TimeLimitSecs, &q.MaxAttempts, &q.ScheduledAt, &q.AvailableUntil, &q.CreatedAt, &q.UpdatedAt, &q.QuestionCount, &q.Weight, &q.VisibleGroupID, &q.UnlockAfterType, &q.UnlockAfterID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Quiz, error) {
	row := r.db.QueryRowContext(ctx, "SELECT "+selectQuizColumns+" FROM quizzes WHERE id = $1", id)
	return scanQuiz(row)
}
