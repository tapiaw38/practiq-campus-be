package assignment

import (
	"context"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	CreateUsecase interface {
		Execute(context.Context, string, bool, string, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		SectionID   *string
		Title       string
		Description string
		DueAt       *time.Time
		MaxScore    int
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

	maxScore := input.MaxScore
	if maxScore <= 0 {
		maxScore = 100
	}

	id, err := app.Repositories.Assignment.Create(ctx, domain.Assignment{
		CourseID:    courseID,
		SectionID:   input.SectionID,
		Title:       input.Title,
		Description: input.Description,
		DueAt:       input.DueAt,
		MaxScore:    maxScore,
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

	return &CreateOutput{Data: toAssignmentData(*created)}, nil
}
