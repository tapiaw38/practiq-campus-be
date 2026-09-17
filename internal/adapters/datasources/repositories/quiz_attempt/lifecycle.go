package quiz_attempt

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

// ErrAttemptsExhausted means the student has already used every attempt the
// quiz allows. Raised inside the start transaction rather than checked before
// it, so two simultaneous requests cannot both pass the check.
var ErrAttemptsExhausted = errors.New("no attempts left for this quiz")

// ErrAttemptClosed means the attempt is already submitted or past its
// deadline, so it can no longer receive answers.
var ErrAttemptClosed = errors.New("attempt is no longer open")

// ErrNotYourAttempt means the attempt belongs to a different student.
var ErrNotYourAttempt = errors.New("attempt belongs to another user")

type (
	// StartParams is what the quiz permits, resolved by the use case before
	// the transaction opens.
	StartParams struct {
		QuizID      string
		UserID      string
		MaxAttempts int
		// TimeLimitSecs is the quiz's allowance per attempt; nil is untimed.
		TimeLimitSecs *int
		// HardDeadline is the quiz's own closing time. It caps the attempt
		// deadline: a 60-minute allowance started 10 minutes before the quiz
		// closes ends when the quiz closes.
		HardDeadline *time.Time
	}

	// StartResult reports which attempt the student should be working on and
	// whether it already existed.
	StartResult struct {
		Attempt domain.QuizAttempt
		Resumed bool
	}
)

// deadlineFor resolves the instant an attempt started at `now` must end by.
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

// StartOrResume hands back the attempt the student should be answering.
//
// Starting used to COUNT attempts and then INSERT count+1, which lost races
// two ways: two simultaneous requests read the same count and collided on
// UNIQUE(quiz_id,user_id,attempt_number), surfacing as a 500; and merely
// reloading the page burned an attempt, because a new row was created while
// the half-finished one was left behind unreachable — on a one-attempt quiz
// that locked the student out of an exam they had not taken.
//
// Both follow from the same thing: starting was not atomic and not
// idempotent. Here an advisory lock serialises the pair (quiz, user) for the
// duration of the transaction — it is taken even when the student has no rows
// yet, which a row lock could not do — and a still-open attempt is returned
// as-is instead of being replaced.
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

	// Serialise this student's starts on this quiz. Released on commit or
	// rollback, so an abandoned transaction cannot wedge the pair.
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

	const attemptColumns = `id, quiz_id, user_id, attempt_number, started_at, expires_at, submitted_at, score, max_score`

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

// SubmitGraded writes the answers and closes the attempt as one unit.
//
// These were two calls with two transactions. A failure between them left the
// answers stored against an attempt that stayed open, and two concurrent
// submissions could both pass the "not submitted yet" check and then
// interleave, saving one set of answers under the other's grade. The close
// was also unconditional, so a late second request would overwrite a finished
// attempt.
//
// One transaction now: the row is locked, its state re-checked under that
// lock, and the UPDATE additionally requires submitted_at IS NULL so a racing
// writer that slipped through loses at the database rather than silently.
// Returns ErrAttemptClosed when the attempt was already finished or expired,
// which the caller reports rather than treating as success.
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

	const attemptColumns = `id, quiz_id, user_id, attempt_number, started_at, expires_at, submitted_at, score, max_score`

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
	if !locked.Live(time.Now()) {
		return nil, ErrAttemptClosed
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
		return nil, ErrAttemptClosed
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return updated, nil
}
