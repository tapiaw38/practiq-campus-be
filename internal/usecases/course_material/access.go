package course_material

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

// requesterOwnsCourse gates writes: only the course owner (or a superadmin)
// adds or removes material.
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

// requesterCanReadCourse gates reads: the owner, a superadmin, or any enrolled
// student — material is what students come to the course for.
func requesterCanReadCourse(ctx context.Context, app *appcontext.Context, requesterID string, isSuperAdmin bool, courseID string) apperrors.ApplicationError {
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
