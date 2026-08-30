package course

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

// requesterOwnsCourse loads the course and confirms the requester is its
// owner or a superadmin. Returns the course so callers that already need it
// (update) don't fetch it twice.
func requesterOwnsCourse(ctx context.Context, app *appcontext.Context, requesterID string, isSuperAdmin bool, courseID string) (*domain.Course, apperrors.ApplicationError) {
	c, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if c == nil {
		return nil, apperrors.NewApplicationError(mappings.CourseNotFoundError, nil)
	}
	if isSuperAdmin || c.OwnerID == requesterID {
		return c, nil
	}
	return nil, apperrors.NewForbiddenError()
}
