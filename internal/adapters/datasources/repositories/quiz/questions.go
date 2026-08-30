package quiz

import (
	"context"
	"encoding/json"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) ListQuestions(ctx context.Context, quizID string) ([]domain.QuizQuestion, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,quiz_id,type,statement,options::text,correct_answer,points,position FROM quiz_questions WHERE quiz_id=$1 ORDER BY position`, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questions := make([]domain.QuizQuestion, 0)
	for rows.Next() {
		var q domain.QuizQuestion
		var optionsJSON string
		if err := rows.Scan(&q.ID, &q.QuizID, &q.Type, &q.Statement, &optionsJSON, &q.CorrectAnswer, &q.Points, &q.Position); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(optionsJSON), &q.Options)
		questions = append(questions, q)
	}
	return questions, rows.Err()
}

func (r *repository) ReplaceQuestions(ctx context.Context, quizID string, in []domain.QuizQuestion) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, `DELETE FROM quiz_questions WHERE quiz_id=$1`, quizID); err != nil {
		return err
	}
	for i, q := range in {
		options := q.Options
		if options == nil {
			options = []string{}
		}
		optionsJSON, err := json.Marshal(options)
		if err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO quiz_questions (quiz_id, type, statement, options, correct_answer, points, position)
			VALUES ($1, $2, $3, $4::jsonb, $5, $6, $7)
		`, quizID, q.Type, q.Statement, string(optionsJSON), q.CorrectAnswer, q.Points, i); err != nil {
			return err
		}
	}
	return tx.Commit()
}
