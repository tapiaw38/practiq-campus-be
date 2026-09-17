package campusaccess

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

// CourseReader is who may see the inside of a course: its staff, or a student
// enrolled in it.
type CourseReader struct {
	// IsManager is true for the owning teacher, the institution's admin, and a
	// superadmin. They see every item, unfiltered and unlocked.
	IsManager bool
}

// EnsureCanReadCourse authorises reading a course's contents.
//
// Being inside the institution is not enough. The tenant-scoped queries decide
// which school a request may touch, and that was mistaken for an answer to
// "may this person open this course?" — so any member could list another
// course's assignments, sections and quizzes by knowing its id. Membership of
// the school and membership of the course are different questions; this one
// answers the second.
//
// Returns whether the caller reads as staff, so callers can skip the group and
// prerequisite filtering that only applies to students.
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
