package quiz

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.Quiz, error) {
	query := `SELECT ` + selectQuizColumns + ` FROM quizzes WHERE course_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	quizzes := make([]domain.Quiz, 0)
	for rows.Next() {
		var q domain.Quiz
		if err := rows.Scan(&q.ID, &q.CourseID, &q.SectionID, &q.Title, &q.Description, &q.TimeLimitSecs, &q.MaxAttempts, &q.ScheduledAt, &q.AvailableUntil, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, q)
	}
	return quizzes, rows.Err()
}
