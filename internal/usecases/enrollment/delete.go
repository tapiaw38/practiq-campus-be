package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	DeleteUsecase interface {
		Execute(context.Context, string, bool, string) apperrors.ApplicationError
	}

	deleteUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{contextFactory: contextFactory}
}

func (u *deleteUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, enrollmentID string) apperrors.ApplicationError {
	app := u.contextFactory()

	e, err := app.Repositories.Enrollment.Get(ctx, enrollmentID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.EnrollmentGetError, err)
	}
	if e == nil {
		return apperrors.NewApplicationError(mappings.EnrollmentNotFoundError, nil)
	}

	if appErr := requesterOwnsCourse(ctx, app, requesterID, isSuperAdmin, e.CourseID); appErr != nil {
		return appErr
	}

	if err := app.Repositories.Enrollment.Delete(ctx, enrollmentID); err != nil {
		return apperrors.NewApplicationError(mappings.EnrollmentDeleteError, err)
	}
	return nil
}
