package quiz

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
		SectionID      *string
		Title          string
		Description    string
		TimeLimitSecs  *int
		MaxAttempts    int
		ScheduledAt    *time.Time
		AvailableUntil *time.Time
	}

	CreateOutput struct {
		Data QuizData `json:"data"`
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

	maxAttempts := input.MaxAttempts
	if maxAttempts < 0 {
		maxAttempts = 1
	}

	id, err := app.Repositories.Quiz.Create(ctx, domain.Quiz{
		CourseID:       courseID,
		SectionID:      input.SectionID,
		Title:          input.Title,
		Description:    input.Description,
		TimeLimitSecs:  input.TimeLimitSecs,
		MaxAttempts:    maxAttempts,
		ScheduledAt:    input.ScheduledAt,
		AvailableUntil: input.AvailableUntil,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.QuizCreateError, err)
	}

	created, err := app.Repositories.Quiz.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.QuizGetError, err)
	}
	if created == nil {
		return nil, apperrors.NewInternalError(nil)
	}

	return &CreateOutput{Data: toQuizData(*created)}, nil
}
