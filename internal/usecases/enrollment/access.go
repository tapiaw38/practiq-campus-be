package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

// requesterOwnsCourse confirms the requester owns the given course or is a
// superadmin. Enrollment intentionally checks course ownership directly
// against the course repository rather than depending on the course usecase
// package, keeping the two features decoupled from each other's internals.
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
