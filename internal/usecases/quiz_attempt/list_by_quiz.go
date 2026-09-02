package quiz_attempt

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/identity"
)

type (
	ListByQuizUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, quizID, bearerToken string) (*ListOutput, apperrors.ApplicationError)
	}

	listByQuizUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewListByQuizUsecase(contextFactory appcontext.Factory) ListByQuizUsecase {
	return &listByQuizUsecase{contextFactory: contextFactory}
}

func (u *listByQuizUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, quizID, bearerToken string) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	q, err := app.Repositories.Quiz.Get(ctx, quizID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if q == nil {
		return nil, apperrors.NewApplicationError(mappings.QuizNotFoundError, nil)
	}
	if !isSuperAdmin {
		course, err := app.Repositories.Course.Get(ctx, q.CourseID)
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		if course == nil || course.OwnerID != requesterID {
			return nil, apperrors.NewForbiddenError()
		}
	}

	attempts, err := app.Repositories.QuizAttempt.ListByQuiz(ctx, quizID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}

	ids := make([]string, 0, len(attempts))
	for _, a := range attempts {
		ids = append(ids, a.UserID)
	}
	names, err := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, ids)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	for i, a := range attempts {
		attempts[i].UserName = identity.FullName(names[a.UserID], a.UserID)
	}

	return &ListOutput{Data: toAttemptDataList(attempts)}, nil
}
