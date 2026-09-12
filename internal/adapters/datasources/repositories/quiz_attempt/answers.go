package quiz_attempt

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// attemptGuard constrains a write to an attempt of the tenant in context,
// walking attempt → quiz → course.
func attemptGuard(ctx context.Context, column string, placeholder int) (string, any, error) {
	return tenantcontext.ExistsIn(ctx,
		"quiz_attempts at JOIN quizzes q ON q.id = at.quiz_id JOIN courses c ON c.id = q.course_id",
		"at.id = "+column, "c.tenant_id", placeholder)
}

// SaveAnswers checks the attempt once before writing, rather than guarding
// each row: every answer in the batch belongs to that one attempt.
func (r *repository) SaveAnswers(ctx context.Context, attemptID string, answers []domain.QuizAnswer) error {
	guard, tenantID, err := attemptGuard(ctx, "$1", 2)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var belongs bool
	if err = tx.QueryRowContext(ctx, `SELECT `+guard, attemptID, tenantID).Scan(&belongs); err != nil {
		return err
	}
	if !belongs {
		return ErrQuizNotInTenant
	}

	for _, a := range answers {
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO quiz_answers (attempt_id, question_id, answer_text, is_correct)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (attempt_id, question_id) DO UPDATE SET answer_text=EXCLUDED.answer_text, is_correct=EXCLUDED.is_correct
		`, attemptID, a.QuestionID, a.AnswerText, a.IsCorrect); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *repository) ListAnswers(ctx context.Context, attemptID string) ([]domain.QuizAnswer, error) {
	tenantID := tenantcontext.ID(ctx)
	if tenantID == "" {
		return nil, tenantcontext.ErrNoTenant
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT ans.id, ans.attempt_id, ans.question_id, ans.answer_text, ans.is_correct
		FROM quiz_answers ans
		JOIN quiz_attempts at ON at.id = ans.attempt_id
		JOIN quizzes q ON q.id = at.quiz_id
		JOIN courses c ON c.id = q.course_id AND c.tenant_id = $2
		WHERE ans.attempt_id = $1`, attemptID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	answers := make([]domain.QuizAnswer, 0)
	for rows.Next() {
		var a domain.QuizAnswer
		if err := rows.Scan(&a.ID, &a.AttemptID, &a.QuestionID, &a.AnswerText, &a.IsCorrect); err != nil {
			return nil, err
		}
		answers = append(answers, a)
	}
	return answers, rows.Err()
}
