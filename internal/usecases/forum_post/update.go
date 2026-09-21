package forum_post

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type UpdateUsecase interface {
	Execute(context.Context, string, string, UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
}

type UpdateInput struct{ Body string }
type UpdateOutput struct {
	Data PostData `json:"data"`
}
type updateUsecase struct{ contextFactory appcontext.Factory }

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, requesterID, postID string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	body := strings.TrimSpace(input.Body)
	if body == "" {
		return nil, apperrors.NewBadRequestError("body is required")
	}

	app := u.contextFactory()
	post, err := app.Repositories.ForumPost.Get(ctx, postID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PostListError, err)
	}
	if post == nil {
		return nil, apperrors.NewNotFoundError("post not found")
	}
	if _, appErr := requesterCanAccessThread(ctx, app, requesterID, false, post.ThreadID); appErr != nil {
		return nil, appErr
	}
	if post.AuthorID != requesterID {
		return nil, apperrors.NewForbiddenError()
	}
	if post.DeletedAt != nil {
		return nil, apperrors.NewBadRequestError("deleted post cannot be updated")
	}
	if err := app.Repositories.ForumPost.Update(ctx, postID, body); err != nil {
		return nil, apperrors.NewApplicationError(mappings.PostCreateError, err)
	}
	post.Body = body
	return &UpdateOutput{Data: toPostData(*post)}, nil
}
