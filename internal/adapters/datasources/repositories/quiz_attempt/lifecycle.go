package quiz_attempt

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

// ErrAttemptsExhausted means the student has already used every attempt the
// quiz allows. Raised inside the start transaction rather than checked before
// it, so two simultaneous requests cannot both pass the check.
var ErrAttemptsExhausted = errors.New("no attempts left for this quiz")

// ErrAttemptAlreadySubmitted means the attempt was handed in already. Kept
// distinct from expiry: telling a student whose time ran out that they
// already submitted describes a submission that never happened, and sends
// them looking for a grade that does not exist.
var ErrAttemptAlreadySubmitted = errors.New("attempt was already submitted")

// ErrAttemptExpired means the deadline passed before the answers arrived.
var ErrAttemptExpired = errors.New("attempt deadline has passed")

// ErrNotYourAttempt means the attempt belongs to a different student.
var ErrNotYourAttempt = errors.New("attempt belongs to another user")

// ErrQuizWindowClosed means the quiz is not open for *starting* right now.
// It never blocks resuming: an attempt already running carries its own
// deadline, fixed when it started.
var ErrQuizWindowClosed = errors.New("quiz is not open for new attempts")

// ErrForeignQuestion means the batch named a question that is not part of the
// attempt's quiz, or named the same one twice.
var ErrForeignQuestion = errors.New("answers reference questions outside this quiz")

// attemptColumns is the row as the domain expects it, in scan order.
const attemptColumns = `id, quiz_id, user_id, attempt_number, started_at, expires_at, submitted_at, score, max_score`

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
		// WindowOpen reports whether the quiz currently admits a *new*
		// attempt. Resuming ignores it: a student halfway through an exam
		// whose teacher moves the dates should finish on the deadline their
		// attempt was given, not be shut out of work already in progress.
		WindowOpen bool
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

	// Nothing to resume, so this would be a new attempt — which the quiz's
	// current window governs. Checked here rather than before the
	// transaction so that a closed window cannot block the resume path.
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

// SaveDraft stores answers without grading or closing the attempt.
//
// Answers used to exist only in the browser's memory until the moment of
// submission, so continuing an attempt showed an empty form — the attempt was
// preserved but everything written into it was gone — and anything that ended
// the page before submitting (a reload, a closed tab, a flat battery) threw
// the work away. They are written to the same rows the submission will
// upsert, so handing in later simply overwrites them.
//
// Refuses once the attempt is submitted or expired: a draft saved after the
// deadline would be work the student was no longer entitled to do.
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
		// is_correct stays false here; grading happens once, at submission.
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

// assertQuestionsBelongTo refuses a batch naming any question that is not part
// of this quiz.
//
// The foreign key only proves a question exists somewhere; it says nothing
// about which quiz, or which institution. Without this, a student could post
// a draft citing a question id belonging to another quiz — another tenant's,
// even — and have the row accepted, attaching an answer to a question the
// attempt has nothing to do with. Checked as a set so a bad batch writes
// nothing rather than partially.
func assertQuestionsBelongTo(ctx context.Context, tx *sql.Tx, quizID string, answers []domain.QuizAnswer) error {
	if len(answers) == 0 {
		return nil
	}
	ids := make([]string, 0, len(answers))
	seen := make(map[string]bool, len(answers))
	for _, a := range answers {
		if seen[a.QuestionID] {
			// The upsert would apply both rows in order, making the result
			// depend on batch ordering. Reject instead of guessing.
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

// FinalizeExpired closes an attempt whose deadline passed without a
// submission, recording the grade for whatever was stored.
//
// Drafts made the work survive, but nothing ever settled those attempts: they
// stayed open forever, and the review screen hides answers until an attempt
// is submitted, so a student who ran out of time had their work saved and
// invisible — to them and to the teacher. Time running out is an ending, and
// this is what makes it one.
//
// Idempotent, and safe to call from any read: the UPDATE only matches an
// attempt that is still open and genuinely past its deadline, so racing
// callers produce one closure and a live attempt is never touched. Returns
// the current row even when another caller already finalized it.
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
	// Someone else finalized it, or it is not actually over: either way there
	// is nothing to do, and saying so is not an error.
	if locked == nil || locked.SubmittedAt != nil || !locked.Expired(time.Now()) {
		return locked, nil
	}

	// Read the draft only after locking the attempt. A save that started
	// before expiry may still be committing while finalization waits here.
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

	// Grades what the student managed to save. Answers they never got to are
	// simply absent, which is the same as leaving them blank.
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
	// Both states are decided here, under the lock, and reported apart: the
	// student needs to know whether their work was recorded or the clock beat
	// them to it.
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
	// The row passed the check above but no longer matches submitted_at IS
	// NULL, so another writer closed it between the two — the only remaining
	// explanation, since this transaction holds the lock.
	if updated == nil {
		return nil, ErrAttemptAlreadySubmitted
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return updated, nil
}
