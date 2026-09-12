package quiz

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// ErrQuizNotInTenant is returned when the quiz belongs to another institution.
var ErrQuizNotInTenant = errors.New("quiz does not belong to this tenant")

func (r *repository) ListQuestions(ctx context.Context, quizID string) ([]domain.QuizQuestion, error) {
	tenantID := tenantcontext.ID(ctx)
	if tenantID == "" {
		return nil, tenantcontext.ErrNoTenant
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT qq.id,qq.quiz_id,qq.type,qq.statement,qq.options::text,qq.correct_answer,qq.points,qq.position
		FROM quiz_questions qq
		JOIN quizzes q ON q.id = qq.quiz_id
		JOIN courses c ON c.id = q.course_id AND c.tenant_id = $2
		WHERE qq.quiz_id=$1 ORDER BY qq.position`, quizID, tenantID)
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

// ReplaceQuestions empties the quiz's question bank before rewriting it, so
// the tenant guard has to hold before the delete: without it a quiz id from
// another institution would wipe that institution's questions.
func (r *repository) ReplaceQuestions(ctx context.Context, quizID string, in []domain.QuizQuestion) error {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"quizzes q JOIN courses c ON c.id = q.course_id", "q.id = $1", "c.tenant_id", 2)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var belongs bool
	if err = tx.QueryRowContext(ctx, `SELECT `+guard, quizID, tenantID).Scan(&belongs); err != nil {
		return err
	}
	if !belongs {
		return ErrQuizNotInTenant
	}

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
