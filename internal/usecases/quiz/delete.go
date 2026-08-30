package quiz

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type DeleteUsecase interface {
	Execute(context.Context, string, bool, string) apperrors.ApplicationError
}

type deleteUsecase struct {
	contextFactory appcontext.Factory
}

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{contextFactory: contextFactory}
}

func (u *deleteUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string) apperrors.ApplicationError {
	app := u.contextFactory()

	existing, err := app.Repositories.Quiz.Get(ctx, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.QuizGetError, err)
	}
	if existing == nil {
		return apperrors.NewApplicationError(mappings.QuizNotFoundError, nil)
	}
	if appErr := requesterOwnsCourse(ctx, app, requesterID, isSuperAdmin, existing.CourseID); appErr != nil {
		return appErr
	}

	if err = app.Repositories.Quiz.Delete(ctx, id); err != nil {
		return apperrors.NewApplicationError(mappings.QuizGetError, err)
	}
	return nil
}
