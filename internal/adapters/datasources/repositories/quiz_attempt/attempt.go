package quiz_attempt

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// ErrQuizNotInTenant is returned when the quiz belongs to another institution.
var ErrQuizNotInTenant = errors.New("quiz does not belong to this tenant")

// Qualified for the tenant chain: attempts reach an institution through their
// quiz's course.
const selectQualifiedAttemptColumns = `
	at.id, at.quiz_id, at.user_id, at.attempt_number, at.started_at, at.submitted_at, at.score, at.max_score
`

// tenantChain is the path every attempt query walks: attempt → quiz → course.
func tenantChain(ctx context.Context) *tenantcontext.Query {
	return tenantcontext.NewQuery(ctx).
		Join("quizzes", "q", "at.quiz_id").
		Through("courses", "c", "q.course_id")
}

// quizGuard constrains a write to a quiz of the tenant in context.
func quizGuard(ctx context.Context, column string, placeholder int) (string, any, error) {
	return tenantcontext.ExistsIn(ctx,
		"quizzes q JOIN courses c ON c.id = q.course_id", "q.id = "+column, "c.tenant_id", placeholder)
}

func scanAttempt(row interface{ Scan(...any) error }) (*domain.QuizAttempt, error) {
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
	guard, tenantID, err := quizGuard(ctx, "$1", 4)
	if err != nil {
		return "", err
	}
	var id string
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO quiz_attempts (quiz_id, user_id, attempt_number)
		SELECT $1, $2, $3 WHERE `+guard+`
		RETURNING id`, a.QuizID, a.UserID, a.AttemptNumber, tenantID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrQuizNotInTenant
	}
	return id, err
}

// Get carries a student's score, so an attempt id from another institution
// must not resolve.
func (r *repository) Get(ctx context.Context, id string) (*domain.QuizAttempt, error) {
	query, args, err := tenantChain(ctx).Where("at.id = ?", id).
		SQL(selectQualifiedAttemptColumns, "quiz_attempts at")
	if err != nil {
		return nil, err
	}
	return scanAttempt(r.db.QueryRowContext(ctx, query, args...))
}

// CountByUser decides whether another attempt is allowed, so it counts within
// the institution the quiz belongs to.
func (r *repository) CountByUser(ctx context.Context, quizID, userID string) (int, error) {
	query, args, err := tenantChain(ctx).
		Where("at.quiz_id = ?", quizID).
		Where("at.user_id = ?", userID).
		SQL("COUNT(*)", "quiz_attempts at")
	if err != nil {
		return 0, err
	}
	var count int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func scanAttempts(rows *sql.Rows) ([]domain.QuizAttempt, error) {
	defer rows.Close()
	attempts := make([]domain.QuizAttempt, 0)
	for rows.Next() {
		a, err := scanAttempt(rows)
		if err != nil {
			return nil, err
		}
		attempts = append(attempts, *a)
	}
	return attempts, rows.Err()
}

func (r *repository) ListByQuiz(ctx context.Context, quizID string) ([]domain.QuizAttempt, error) {
	query, args, err := tenantChain(ctx).Where("at.quiz_id = ?", quizID).
		SQL(selectQualifiedAttemptColumns, "quiz_attempts at")
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query+" ORDER BY at.started_at DESC", args...)
	if err != nil {
		return nil, err
	}
	return scanAttempts(rows)
}

func (r *repository) ListMine(ctx context.Context, quizID, userID string) ([]domain.QuizAttempt, error) {
	query, args, err := tenantChain(ctx).
		Where("at.quiz_id = ?", quizID).
		Where("at.user_id = ?", userID).
		SQL(selectQualifiedAttemptColumns, "quiz_attempts at")
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query+" ORDER BY at.attempt_number ASC", args...)
	if err != nil {
		return nil, err
	}
	return scanAttempts(rows)
}

func (r *repository) Submit(ctx context.Context, id string, score, maxScore int) error {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"quizzes q JOIN courses c ON c.id = q.course_id",
		"q.id = quiz_attempts.quiz_id", "c.tenant_id", 4)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx,
		`UPDATE quiz_attempts SET submitted_at=NOW(), score=$1, max_score=$2 WHERE id=$3 AND `+guard,
		score, maxScore, id, tenantID)
	return err
}
