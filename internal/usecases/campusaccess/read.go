package campusaccess

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type CourseReader struct {
	IsManager bool
}

func EnsureCanReadCourse(ctx context.Context, app *appcontext.Context, requesterID string, isSuperAdmin bool, courseID string) (CourseReader, apperrors.ApplicationError) {
	if isSuperAdmin {
		return CourseReader{IsManager: true}, nil
	}

	course, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return CourseReader{}, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if course == nil {
		return CourseReader{}, apperrors.NewApplicationError(mappings.CourseNotFoundError, nil)
	}
	if CanManageCourse(ctx, course.OwnerID, requesterID, false) {
		return CourseReader{IsManager: true}, nil
	}

	enrollment, err := app.Repositories.Enrollment.GetByCourseAndUser(ctx, courseID, requesterID)
	if err != nil {
		return CourseReader{}, apperrors.NewApplicationError(mappings.EnrollmentGetError, err)
	}
	if enrollment == nil {
		return CourseReader{}, apperrors.NewForbiddenError()
	}
	return CourseReader{}, nil
}
