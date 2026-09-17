package quiz_attempt

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
)

// finalizeIfExpired settles an attempt whose time ran out without a hand-in.
//
// Saving drafts made a student's work survive the clock, but nothing ever
// ended those attempts: they stayed open indefinitely, and every screen that
// shows answers waits for submitted_at, so the work sat in the database
// invisible to the student and to the teacher. Running out of time is an
// ending; this records it, grading exactly what was saved and leaving
// unanswered questions unanswered.
//
// Called from the read paths rather than a scheduled job: the moment anyone
// looks at an attempt is a perfectly good moment to notice it is over, and it
// needs no process that can be forgotten or fall behind. The repository makes
// it idempotent and a no-op on anything still live, so concurrent readers
// settle it once between them.
//
// Returns the attempt to display: the finalized one when it just closed, or
// the original unchanged.
func finalizeIfExpired(ctx context.Context, app *appcontext.Context, attempt domain.QuizAttempt) (domain.QuizAttempt, error) {
	if attempt.SubmittedAt != nil || !attempt.Expired(time.Now()) {
		return attempt, nil
	}

	// Read and grade the draft under the repository's attempt lock, so a
	// save still committing at the deadline cannot be omitted from the grade.
	finalized, err := app.Repositories.QuizAttempt.FinalizeExpired(ctx, attempt.ID)
	if err != nil {
		return attempt, err
	}
	// nil means it was not actually over, or somebody else closed it first.
	if finalized == nil {
		return attempt, nil
	}
	return *finalized, nil
}

// finalizeExpiredIn settles every attempt in a list that is past its deadline,
// so a teacher opening the results page sees finished work rather than a
// column of rows stuck on "in progress".
func finalizeExpiredIn(ctx context.Context, app *appcontext.Context, attempts []domain.QuizAttempt) ([]domain.QuizAttempt, error) {
	for i, attempt := range attempts {
		settled, err := finalizeIfExpired(ctx, app, attempt)
		if err != nil {
			return attempts, err
		}
		attempts[i] = settled
	}
	return attempts, nil
}
