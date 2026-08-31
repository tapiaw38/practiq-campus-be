package assignment

import (
	"context"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/unlock"
)

type (
	CreateUsecase interface {
		Execute(context.Context, string, bool, string, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		SectionID       *string
		Title           string
		Description     string
		DueAt           *time.Time
		MaxScore        int
		Weight          int
		VisibleGroupID  *string
		UnlockAfterType *string
		UnlockAfterID   *string
	}

	CreateOutput struct {
		Data AssignmentData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, apperrors.NewBadRequestError("title is required")
	}

	app := u.contextFactory()

	if appErr := requesterOwnsCourse(ctx, app, requesterID, isSuperAdmin, courseID); appErr != nil {
		return nil, appErr
	}
	if appErr := validateVisibleGroup(ctx, app, courseID, input.VisibleGroupID); appErr != nil {
		return nil, appErr
	}
	if appErr := validateUnlockTarget(ctx, app, courseID, "", input.UnlockAfterType, input.UnlockAfterID); appErr != nil {
		return nil, appErr
	}

	maxScore := input.MaxScore
	if maxScore <= 0 {
		maxScore = 100
	}
	weight := input.Weight
	if weight <= 0 {
		weight = 100
	}

	id, err := app.Repositories.Assignment.Create(ctx, domain.Assignment{
		CourseID:        courseID,
		SectionID:       input.SectionID,
		Title:           input.Title,
		Description:     input.Description,
		DueAt:           input.DueAt,
		MaxScore:        maxScore,
		Weight:          weight,
		VisibleGroupID:  input.VisibleGroupID,
		UnlockAfterType: input.UnlockAfterType,
		UnlockAfterID:   input.UnlockAfterID,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentCreateError, err)
	}

	created, err := app.Repositories.Assignment.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentGetError, err)
	}
	if created == nil {
		return nil, apperrors.NewInternalError(nil)
	}
	if enrollments, err := app.Repositories.Enrollment.ListByCourse(ctx, courseID); err == nil {
		for _, enrollment := range enrollments {
			_ = app.Repositories.Notification.Create(ctx, domain.Notification{UserID: enrollment.UserID, Type: "assignment_created", Title: "Nueva tarea: " + created.Title, Body: "Tenés una nueva tarea", Data: `{"assignment_id":"` + created.ID + `","course_id":"` + courseID + `"}`})
		}
	}

	return &CreateOutput{Data: toAssignmentData(*created, unlock.Status{})}, nil
}
