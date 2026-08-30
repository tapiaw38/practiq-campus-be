package quiz_attempt

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

const selectAttemptColumns = `
	id, quiz_id, user_id, attempt_number, started_at, submitted_at, score, max_score
`

func scanAttempt(row *sql.Row) (*domain.QuizAttempt, error) {
	var a domain.QuizAttempt
	err := row.Scan(&a.ID, &a.QuizID, &a.UserID, &a.AttemptNumber, &a.StartedAt, &a.SubmittedAt, &a.Score, &a.MaxScore)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repository) Create(ctx context.Context, a domain.QuizAttempt) (string, error) {
	query := `
		INSERT INTO quiz_attempts (quiz_id, user_id, attempt_number)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, a.QuizID, a.UserID, a.AttemptNumber).Scan(&id)
	return id, err
}

func (r *repository) Get(ctx context.Context, id string) (*domain.QuizAttempt, error) {
	row := r.db.QueryRowContext(ctx, "SELECT "+selectAttemptColumns+" FROM quiz_attempts WHERE id = $1", id)
	return scanAttempt(row)
}

func (r *repository) CountByUser(ctx context.Context, quizID, userID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM quiz_attempts WHERE quiz_id=$1 AND user_id=$2`, quizID, userID).Scan(&count)
	return count, err
}

func scanAttempts(rows *sql.Rows) ([]domain.QuizAttempt, error) {
	defer rows.Close()
	attempts := make([]domain.QuizAttempt, 0)
	for rows.Next() {
		var a domain.QuizAttempt
		if err := rows.Scan(&a.ID, &a.QuizID, &a.UserID, &a.AttemptNumber, &a.StartedAt, &a.SubmittedAt, &a.Score, &a.MaxScore); err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}
	return attempts, rows.Err()
}

func (r *repository) ListByQuiz(ctx context.Context, quizID string) ([]domain.QuizAttempt, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT a.id, a.quiz_id, a.user_id, COALESCE(p.full_name, a.user_id), a.attempt_number, a.started_at, a.submitted_at, a.score, a.max_score
		FROM quiz_attempts a
		LEFT JOIN campus_profiles p ON p.id = a.user_id
		WHERE a.quiz_id=$1 ORDER BY a.started_at DESC
	`, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attempts := make([]domain.QuizAttempt, 0)
	for rows.Next() {
		var a domain.QuizAttempt
		if err := rows.Scan(&a.ID, &a.QuizID, &a.UserID, &a.UserName, &a.AttemptNumber, &a.StartedAt, &a.SubmittedAt, &a.Score, &a.MaxScore); err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}
	return attempts, rows.Err()
}

func (r *repository) ListMine(ctx context.Context, quizID, userID string) ([]domain.QuizAttempt, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT "+selectAttemptColumns+" FROM quiz_attempts WHERE quiz_id=$1 AND user_id=$2 ORDER BY attempt_number ASC", quizID, userID)
	if err != nil {
		return nil, err
	}
	return scanAttempts(rows)
}

func (r *repository) Submit(ctx context.Context, id string, score, maxScore int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE quiz_attempts SET submitted_at=NOW(), score=$1, max_score=$2 WHERE id=$3`, score, maxScore, id)
	return err
}
