package quiz

import (
	"context"
	"strings"
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
		SectionID      *string
		Title          string
		Description    string
		TimeLimitSecs  *int
		MaxAttempts    int
		ScheduledAt    *time.Time
		AvailableUntil *time.Time
	}

	UpdateOutput struct {
		Data QuizData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, apperrors.NewBadRequestError("title is required")
	}

	app := u.contextFactory()

	existing, err := app.Repositories.Quiz.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.QuizGetError, err)
	}
	if existing == nil {
		return nil, apperrors.NewApplicationError(mappings.QuizNotFoundError, nil)
	}
	if appErr := requesterOwnsCourse(ctx, app, requesterID, isSuperAdmin, existing.CourseID); appErr != nil {
		return nil, appErr
	}

	maxAttempts := input.MaxAttempts
	if maxAttempts < 0 {
		maxAttempts = 1
	}

	existing.SectionID = input.SectionID
	existing.Title = input.Title
	existing.Description = input.Description
	existing.TimeLimitSecs = input.TimeLimitSecs
	existing.MaxAttempts = maxAttempts
	existing.ScheduledAt = input.ScheduledAt
	existing.AvailableUntil = input.AvailableUntil

	if err = app.Repositories.Quiz.Update(ctx, id, *existing); err != nil {
		return nil, apperrors.NewApplicationError(mappings.QuizGetError, err)
	}

	return &UpdateOutput{Data: toQuizData(*existing)}, nil
}
