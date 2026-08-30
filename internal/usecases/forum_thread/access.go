package forum_thread

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

// requesterCanAccessCourse is deliberately looser than requesterOwnsCourse
// elsewhere: a forum is for the course owner AND every enrolled student, not
// just the teacher — posting is the whole point of a forum.
func requesterCanAccessCourse(ctx context.Context, app *appcontext.Context, requesterID string, isSuperAdmin bool, courseID string) apperrors.ApplicationError {
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
	if c.OwnerID == requesterID {
		return nil
	}

	enrollment, err := app.Repositories.Enrollment.GetByCourseAndUser(ctx, courseID, requesterID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.EnrollmentGetError, err)
	}
	if enrollment == nil {
		return apperrors.NewForbiddenError()
	}
	return nil
}

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
