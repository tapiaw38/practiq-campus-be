package course_group

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
)

func (u *usecase) Create(ctx context.Context, requesterID string, isSuperAdmin bool, courseID, name string) (*GroupData, apperrors.ApplicationError) {
	if strings.TrimSpace(name) == "" {
		return nil, apperrors.NewBadRequestError("name is required")
	}
	if appErr := u.requesterOwnsGroupsCourse(ctx, requesterID, isSuperAdmin, courseID); appErr != nil {
		return nil, appErr
	}
	app := u.f()
	id, err := app.Repositories.CourseGroup.Create(ctx, domain.CourseGroup{CourseID: courseID, Name: strings.TrimSpace(name)})
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	created, err := app.Repositories.CourseGroup.Get(ctx, id)
	if err != nil || created == nil {
		return nil, apperrors.NewInternalError(err)
	}
	data := toGroupData(*created)
	return &data, nil
}
