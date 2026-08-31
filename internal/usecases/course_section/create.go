package course_section

import (
	"context"
	"strings"

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
		Title       string
		Description string
	}

	CreateOutput struct {
		Data SectionData `json:"data"`
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

	existing, err := app.Repositories.CourseSection.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SectionListError, err)
	}

	id, err := app.Repositories.CourseSection.Create(ctx, domain.CourseSection{
		CourseID:    courseID,
		Title:       input.Title,
		Description: input.Description,
		Position:    len(existing) + 1,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SectionCreateError, err)
	}

	created, err := app.Repositories.CourseSection.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SectionGetError, err)
	}
	if created == nil {
		return nil, apperrors.NewInternalError(nil)
	}

	return &CreateOutput{Data: toSectionData(*created)}, nil
}
