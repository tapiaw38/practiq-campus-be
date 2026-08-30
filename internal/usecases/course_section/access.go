package course_section

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

// requesterOwnsCourse checks course ownership directly against the course
// repository rather than depending on the course usecase package — same
// decoupling choice already made for enrollment.
func requesterOwnsCourse(ctx context.Context, app *appcontext.Context, requesterID string, isSuperAdmin bool, courseID string) apperrors.ApplicationError {
	if isSuperAdmin {
		return nil
	}
	c, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if c == nil {
		return apperrors.NewApplicationError(mappings.CourseNotFoundError, nil)
	}
	if c.OwnerID != requesterID {
		return apperrors.NewForbiddenError()
	}
	return nil
}
