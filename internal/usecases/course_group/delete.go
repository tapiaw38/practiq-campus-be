package course_group

import (
	"context"

	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
)

func (u *usecase) Delete(ctx context.Context, requesterID string, isSuperAdmin bool, groupID string) apperrors.ApplicationError {
	courseID, appErr := u.groupCourseID(ctx, groupID)
	if appErr != nil {
		return appErr
	}
	if appErr := u.requesterOwnsGroupsCourse(ctx, requesterID, isSuperAdmin, courseID); appErr != nil {
		return appErr
	}
	if err := u.f().Repositories.CourseGroup.Delete(ctx, groupID); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}
