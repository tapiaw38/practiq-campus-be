package quiz_attempt

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
)

func finalizeIfExpired(ctx context.Context, app *appcontext.Context, attempt domain.QuizAttempt) (domain.QuizAttempt, error) {
	if attempt.SubmittedAt != nil || !attempt.Expired(time.Now()) {
		return attempt, nil
	}

	finalized, err := app.Repositories.QuizAttempt.FinalizeExpired(ctx, attempt.ID)
	if err != nil {
		return attempt, err
	}

	if finalized == nil {
		return attempt, nil
	}
	return *finalized, nil
}

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
