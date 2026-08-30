package course_group

import (
	"context"

	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
)

func (u *usecase) AddMember(ctx context.Context, requesterID string, isSuperAdmin bool, groupID, userID string) apperrors.ApplicationError {
	courseID, appErr := u.groupCourseID(ctx, groupID)
	if appErr != nil {
		return appErr
	}
	if appErr := u.requesterOwnsGroupsCourse(ctx, requesterID, isSuperAdmin, courseID); appErr != nil {
		return appErr
	}
	app := u.f()
	enrollment, err := app.Repositories.Enrollment.GetByCourseAndUser(ctx, courseID, userID)
	if err != nil {
		return apperrors.NewInternalError(err)
	}
	if enrollment == nil {
		return apperrors.NewBadRequestError("user is not enrolled in this course")
	}
	if err := app.Repositories.CourseGroup.AddMember(ctx, groupID, userID); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (u *usecase) RemoveMember(ctx context.Context, requesterID string, isSuperAdmin bool, groupID, userID string) apperrors.ApplicationError {
	courseID, appErr := u.groupCourseID(ctx, groupID)
	if appErr != nil {
		return appErr
	}
	if appErr := u.requesterOwnsGroupsCourse(ctx, requesterID, isSuperAdmin, courseID); appErr != nil {
		return appErr
	}
	if err := u.f().Repositories.CourseGroup.RemoveMember(ctx, groupID, userID); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}
