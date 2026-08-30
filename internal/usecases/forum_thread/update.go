package forum_thread

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type UpdateUsecase interface {
	Execute(context.Context, string, bool, string, UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
}

type UpdateInput struct{ Title, Description string }
type UpdateOutput struct {
	Data ThreadData `json:"data"`
}
type updateUsecase struct{ contextFactory appcontext.Factory }

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, threadID string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, apperrors.NewBadRequestError("title is required")
	}
	app := u.contextFactory()
	thread, err := app.Repositories.ForumThread.Get(ctx, threadID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ThreadGetError, err)
	}
	if thread == nil {
		return nil, apperrors.NewApplicationError(mappings.ThreadNotFoundError, nil)
	}
	if appErr := requesterOwnsCourse(ctx, app, requesterID, isSuperAdmin, thread.CourseID); appErr != nil {
		return nil, appErr
	}
	if err := app.Repositories.ForumThread.Update(ctx, threadID, domain.ForumThread{Title: strings.TrimSpace(input.Title), Description: strings.TrimSpace(input.Description)}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ThreadUpdateError, err)
	}
	updated, err := app.Repositories.ForumThread.Get(ctx, threadID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ThreadGetError, err)
	}
	if updated == nil {
		return nil, apperrors.NewInternalError(nil)
	}
	return &UpdateOutput{Data: toThreadData(*updated)}, nil
}
