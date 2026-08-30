package course_material

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

// Execute removes the row but deliberately leaves the S3 object in place:
// deleting is cheap to undo when the row is gone but the bytes are not, and
// an orphaned object costs storage, not correctness.
func (u *deleteUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, materialID string) apperrors.ApplicationError {
	app := u.contextFactory()

	material, err := app.Repositories.CourseMaterial.Get(ctx, materialID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.MaterialGetError, err)
	}
	if material == nil {
		return apperrors.NewApplicationError(mappings.MaterialNotFoundError, nil)
	}

	if appErr := requesterOwnsCourse(ctx, app, requesterID, isSuperAdmin, material.CourseID); appErr != nil {
		return appErr
	}

	if err := app.Repositories.CourseMaterial.Delete(ctx, materialID); err != nil {
		return apperrors.NewApplicationError(mappings.MaterialDeleteError, err)
	}
	return nil
}
