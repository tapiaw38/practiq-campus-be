package quiz_attempt

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) SaveAnswers(ctx context.Context, attemptID string, answers []domain.QuizAnswer) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

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
	rows, err := r.db.QueryContext(ctx, `SELECT id,attempt_id,question_id,answer_text,is_correct FROM quiz_answers WHERE attempt_id=$1`, attemptID)
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
