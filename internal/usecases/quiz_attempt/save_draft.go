package quiz_attempt

import (
	"context"
	"errors"

	attemptrepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/quiz_attempt"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	SaveDraftUsecase interface {
		Execute(ctx context.Context, requesterID, attemptID string, answers []AnswerInput) apperrors.ApplicationError
	}

	saveDraftUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewSaveDraftUsecase(contextFactory appcontext.Factory) SaveDraftUsecase {
	return &saveDraftUsecase{contextFactory: contextFactory}
}

func (u *saveDraftUsecase) Execute(ctx context.Context, requesterID, attemptID string, answers []AnswerInput) apperrors.ApplicationError {
	app := u.contextFactory()

	drafts := make([]domain.QuizAnswer, 0, len(answers))
	for _, a := range answers {
		drafts = append(drafts, domain.QuizAnswer{
			AttemptID:  attemptID,
			QuestionID: a.QuestionID,
			AnswerText: a.AnswerText,
		})
	}

	err := app.Repositories.QuizAttempt.SaveDraft(ctx, attemptID, requesterID, drafts)
	switch {
	case errors.Is(err, attemptrepo.ErrNotYourAttempt):
		return apperrors.NewForbiddenError()
	case errors.Is(err, attemptrepo.ErrAttemptAlreadySubmitted):
		return apperrors.NewApplicationError(mappings.QuizAttemptAlreadySubmittedError, nil)
	case errors.Is(err, attemptrepo.ErrAttemptExpired):
		return apperrors.NewApplicationError(mappings.QuizAttemptExpiredError, nil)
	case err != nil:
		return apperrors.NewInternalError(err)
	}
	return nil
}
