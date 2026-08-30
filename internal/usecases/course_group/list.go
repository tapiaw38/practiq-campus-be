package course_group

import (
	"context"

	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
)

func (u *usecase) List(ctx context.Context, courseID string) ([]GroupData, apperrors.ApplicationError) {
	groups, err := u.f().Repositories.CourseGroup.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return toGroupDataList(groups), nil
}
