package quiz

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.Quiz, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "q.course_id").
		Where("q.course_id = ?", courseID).
		SQL(selectQualifiedQuizColumns, "quizzes q")
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query+" ORDER BY q.created_at ASC", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []domain.Quiz
	for rows.Next() {
		q, err := scanQuiz(rows)
		if err != nil {
			return nil, err
		}
		quizzes = append(quizzes, *q)
	}
	return quizzes, rows.Err()
}
