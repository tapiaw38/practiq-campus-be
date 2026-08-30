package course

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type DeleteUsecase interface {
	Execute(context.Context, string, bool, string) apperrors.ApplicationError
}

type deleteUsecase struct{ contextFactory appcontext.Factory }

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{contextFactory: contextFactory}
}

func (u *deleteUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string) apperrors.ApplicationError {
	app := u.contextFactory()
	if _, appErr := requesterOwnsCourse(ctx, app, requesterID, isSuperAdmin, courseID); appErr != nil {
		return appErr
	}
	if err := app.Repositories.Course.Delete(ctx, courseID); err != nil {
		return apperrors.NewApplicationError(mappings.CourseUpdateError, err)
	}
	return nil
}
