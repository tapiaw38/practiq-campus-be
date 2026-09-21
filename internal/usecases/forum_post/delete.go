package forum_post

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type DeleteUsecase interface {
	Execute(context.Context, string, string) apperrors.ApplicationError
}
type deleteUsecase struct{ contextFactory appcontext.Factory }

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{contextFactory: contextFactory}
}

func (u *deleteUsecase) Execute(ctx context.Context, requesterID, postID string) apperrors.ApplicationError {
	app := u.contextFactory()
	post, err := app.Repositories.ForumPost.Get(ctx, postID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.PostListError, err)
	}
	if post == nil {
		return apperrors.NewNotFoundError("post not found")
	}
	if _, appErr := requesterCanAccessThread(ctx, app, requesterID, false, post.ThreadID); appErr != nil {
		return appErr
	}
	if post.AuthorID != requesterID {
		return apperrors.NewForbiddenError()
	}
	if post.DeletedAt != nil {
		return nil
	}
	if err := app.Repositories.ForumPost.Delete(ctx, postID); err != nil {
		return apperrors.NewApplicationError(mappings.PostCreateError, err)
	}
	return nil
}
