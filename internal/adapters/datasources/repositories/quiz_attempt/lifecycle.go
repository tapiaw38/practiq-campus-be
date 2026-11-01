package quiz_attempt

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

var ErrAttemptsExhausted = errors.New("no attempts left for this quiz")

var ErrAttemptAlreadySubmitted = errors.New("attempt was already submitted")

var ErrAttemptExpired = errors.New("attempt deadline has passed")

var ErrNotYourAttempt = errors.New("attempt belongs to another user")

var ErrQuizWindowClosed = errors.New("quiz is not open for new attempts")

var ErrForeignQuestion = errors.New("answers reference questions outside this quiz")

const attemptColumns = `id, quiz_id, user_id, attempt_number, started_at, expires_at, submitted_at, score, max_score`

type (
	StartParams struct {
		QuizID      string
		UserID      string
		MaxAttempts int

		TimeLimitSecs *int

		HardDeadline *time.Time

		WindowOpen bool
	}

	StartResult struct {
		Attempt domain.QuizAttempt
		Resumed bool
	}
)

func deadlineFor(now time.Time, p StartParams) *time.Time {
	var deadline *time.Time
	if p.TimeLimitSecs != nil && *p.TimeLimitSecs > 0 {
		d := now.Add(time.Duration(*p.TimeLimitSecs) * time.Second)
		deadline = &d
	}
	if p.HardDeadline != nil && (deadline == nil || p.HardDeadline.Before(*deadline)) {
		d := *p.HardDeadline
		deadline = &d
	}
	return deadline
}

func (r *repository) StartOrResume(ctx context.Context, p StartParams) (*StartResult, error) {
	guard, tenantID, err := quizGuard(ctx, "$1", 2)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx,
		`SELECT pg_advisory_xact_lock(hashtext($1 || ':' || $2))`, p.QuizID, p.UserID); err != nil {
		return nil, err
	}

	var belongs bool
	if err = tx.QueryRowContext(ctx, `SELECT `+guard, p.QuizID, tenantID).Scan(&belongs); err != nil {
		return nil, err
	}
	if !belongs {
		return nil, ErrQuizNotInTenant
	}

	live, err := scanAttempt(tx.QueryRowContext(ctx, `
		SELECT `+attemptColumns+`
		FROM quiz_attempts
		WHERE quiz_id = $1 AND user_id = $2
		  AND submitted_at IS NULL
		  AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY attempt_number DESC
		LIMIT 1`, p.QuizID, p.UserID))
	if err != nil {
		return nil, err
	}
	if live != nil {
		if err = tx.Commit(); err != nil {
			return nil, err
		}
		return &StartResult{Attempt: *live, Resumed: true}, nil
	}

	if !p.WindowOpen {
		return nil, ErrQuizWindowClosed
	}

	var used, highest int
	if err = tx.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(MAX(attempt_number), 0)
		FROM quiz_attempts WHERE quiz_id = $1 AND user_id = $2`,
		p.QuizID, p.UserID).Scan(&used, &highest); err != nil {
		return nil, err
	}
	if p.MaxAttempts > 0 && used >= p.MaxAttempts {
		return nil, ErrAttemptsExhausted
	}

	created, err := scanAttempt(tx.QueryRowContext(ctx, `
		INSERT INTO quiz_attempts (quiz_id, user_id, attempt_number, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING `+attemptColumns,
		p.QuizID, p.UserID, highest+1, deadlineFor(time.Now(), p)))
	if err != nil {
		return nil, err
	}
	if created == nil {
		return nil, sql.ErrNoRows
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &StartResult{Attempt: *created}, nil
}

func (r *repository) SaveDraft(ctx context.Context, attemptID, userID string, answers []domain.QuizAnswer) error {
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

	locked, err := scanAttempt(tx.QueryRowContext(ctx,
		`SELECT id, quiz_id, user_id, attempt_number, started_at, expires_at, submitted_at, score, max_score
		 FROM quiz_attempts WHERE id = $1 FOR UPDATE`, attemptID))
	if err != nil {
		return err
	}
	if locked == nil {
		return sql.ErrNoRows
	}
	if locked.UserID != userID {
		return ErrNotYourAttempt
	}
	if locked.SubmittedAt != nil {
		return ErrAttemptAlreadySubmitted
	}
	if locked.Expired(time.Now()) {
		return ErrAttemptExpired
	}

	if err = assertQuestionsBelongTo(ctx, tx, locked.QuizID, answers); err != nil {
		return err
	}

	for _, a := range answers {

		if _, err = tx.ExecContext(ctx, `
			INSERT INTO quiz_answers (attempt_id, question_id, answer_text, is_correct)
			VALUES ($1, $2, $3, FALSE)
			ON CONFLICT (attempt_id, question_id)
			DO UPDATE SET answer_text = EXCLUDED.answer_text
		`, attemptID, a.QuestionID, a.AnswerText); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func assertQuestionsBelongTo(ctx context.Context, tx *sql.Tx, quizID string, answers []domain.QuizAnswer) error {
	if len(answers) == 0 {
		return nil
	}
	ids := make([]string, 0, len(answers))
	seen := make(map[string]bool, len(answers))
	for _, a := range answers {
		if seen[a.QuestionID] {

			return ErrForeignQuestion
		}
		seen[a.QuestionID] = true
		ids = append(ids, a.QuestionID)
	}

	var matching int
	if err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM quiz_questions WHERE quiz_id = $1 AND id = ANY($2)`,
		quizID, pq.Array(ids)).Scan(&matching); err != nil {
		return err
	}
	if matching != len(ids) {
		return ErrForeignQuestion
	}
	return nil
}

func (r *repository) FinalizeExpired(ctx context.Context, attemptID string) (*domain.QuizAttempt, error) {
	guard, tenantID, err := attemptGuard(ctx, "$1", 2)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var belongs bool
	if err = tx.QueryRowContext(ctx, `SELECT `+guard, attemptID, tenantID).Scan(&belongs); err != nil {
		return nil, err
	}
	if !belongs {
		return nil, ErrQuizNotInTenant
	}

	locked, err := scanAttempt(tx.QueryRowContext(ctx,
		`SELECT `+attemptColumns+` FROM quiz_attempts WHERE id = $1 FOR UPDATE`, attemptID))
	if err != nil {
		return nil, err
	}

	if locked == nil || locked.SubmittedAt != nil || !locked.Expired(time.Now()) {
		return locked, nil
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT q.id, q.type, q.statement, q.correct_answer, q.points, COALESCE(a.answer_text, '')
		FROM quiz_questions q
		LEFT JOIN quiz_answers a ON a.question_id = q.id AND a.attempt_id = $1
		WHERE q.quiz_id = $2 ORDER BY q.position`, attemptID, locked.QuizID)
	if err != nil {
		return nil, err
	}
	var answers []domain.QuizAnswer
	score, maxScore := 0, 0
	for rows.Next() {
		var question domain.QuizQuestion
		var answer string
		if err = rows.Scan(&question.ID, &question.Type, &question.Statement, &question.CorrectAnswer, &question.Points, &answer); err != nil {
			rows.Close()
			return nil, err
		}
		correct := domain.GradeQuizAnswer(question, answer)
		maxScore += question.Points
		if correct {
			score += question.Points
		}
		answers = append(answers, domain.QuizAnswer{QuestionID: question.ID, AnswerText: answer, IsCorrect: correct})
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}

	for _, a := range answers {
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO quiz_answers (attempt_id, question_id, answer_text, is_correct)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (attempt_id, question_id) DO UPDATE SET is_correct = EXCLUDED.is_correct
		`, attemptID, a.QuestionID, a.AnswerText, a.IsCorrect); err != nil {
			return nil, err
		}
	}

	updated, err := scanAttempt(tx.QueryRowContext(ctx, `
		UPDATE quiz_attempts SET submitted_at = NOW(), score = $1, max_score = $2
		WHERE id = $3 AND submitted_at IS NULL AND expires_at IS NOT NULL AND expires_at <= NOW()
		RETURNING `+attemptColumns, score, maxScore, attemptID))
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, nil
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *repository) SubmitGraded(ctx context.Context, attemptID, userID string, answers []domain.QuizAnswer, score, maxScore int) (*domain.QuizAttempt, error) {
	guard, tenantID, err := attemptGuard(ctx, "$1", 2)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var belongs bool
	if err = tx.QueryRowContext(ctx, `SELECT `+guard, attemptID, tenantID).Scan(&belongs); err != nil {
		return nil, err
	}
	if !belongs {
		return nil, ErrQuizNotInTenant
	}

	locked, err := scanAttempt(tx.QueryRowContext(ctx,
		`SELECT `+attemptColumns+` FROM quiz_attempts WHERE id = $1 FOR UPDATE`, attemptID))
	if err != nil {
		return nil, err
	}
	if locked == nil {
		return nil, sql.ErrNoRows
	}
	if locked.UserID != userID {
		return nil, ErrNotYourAttempt
	}

	if locked.SubmittedAt != nil {
		return nil, ErrAttemptAlreadySubmitted
	}
	if locked.Expired(time.Now()) {
		return nil, ErrAttemptExpired
	}

	for _, a := range answers {
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO quiz_answers (attempt_id, question_id, answer_text, is_correct)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (attempt_id, question_id)
			DO UPDATE SET answer_text = EXCLUDED.answer_text, is_correct = EXCLUDED.is_correct
		`, attemptID, a.QuestionID, a.AnswerText, a.IsCorrect); err != nil {
			return nil, err
		}
	}

	updated, err := scanAttempt(tx.QueryRowContext(ctx, `
		UPDATE quiz_attempts SET submitted_at = NOW(), score = $1, max_score = $2
		WHERE id = $3 AND submitted_at IS NULL
		RETURNING `+attemptColumns, score, maxScore, attemptID))
	if err != nil {
		return nil, err
	}

	if updated == nil {
		return nil, ErrAttemptAlreadySubmitted
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return updated, nil
}
