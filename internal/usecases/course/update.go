package course

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	UpdateUsecase interface {
		Execute(context.Context, string, bool, string, UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateInput struct {
		Title            string
		Description      string
		Status           string
		StartDate        *time.Time
		EndDate          *time.Time
		PractiqSubjectID *string
		Labels           *[]string
	}

	UpdateOutput struct {
		Data CourseData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

var validCourseStatus = map[string]bool{
	"draft": true, "published": true, "archived": true,
}

func (u *updateUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	current, appErr := requesterOwnsCourse(ctx, app, requesterID, isSuperAdmin, courseID)
	if appErr != nil {
		return nil, appErr
	}

	status := input.Status
	if status == "" {
		status = current.Status
	} else if !validCourseStatus[status] {
		return nil, apperrors.NewBadRequestError("invalid status")
	}

	title := input.Title
	if title == "" {
		title = current.Title
	}

	current.Title = title
	current.Description = input.Description
	current.Status = status
	current.StartDate = input.StartDate
	current.EndDate = input.EndDate
	current.PractiqSubjectID = input.PractiqSubjectID
	if input.Labels != nil {
		current.Labels = normalizeLabels(*input.Labels)
	}

	if err := app.Repositories.Course.Update(ctx, courseID, *current); err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseUpdateError, err)
	}

	updated, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if updated == nil {
		return nil, apperrors.NewInternalError(nil)
	}

	return &UpdateOutput{Data: toCourseData(*updated)}, nil
}
