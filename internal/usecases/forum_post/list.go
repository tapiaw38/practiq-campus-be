package forum_post

import (
	"context"
	repo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/forum_post"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(context.Context, string, bool, string, repo.ListOptions) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data    []PostData `json:"data"`
		HasMore bool       `json:"has_more"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, threadID string, options repo.ListOptions) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if _, appErr := requesterCanAccessThread(ctx, app, requesterID, isSuperAdmin, threadID); appErr != nil {
		return nil, appErr
	}

	posts, err := app.Repositories.ForumPost.ListByThread(ctx, threadID, options)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PostListError, err)
	}

	rootCount, err := app.Repositories.ForumPost.CountRootsByThread(ctx, threadID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PostListError, err)
	}
	limit := options.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return &ListOutput{
		Data:    toPostDataList(posts),
		HasMore: options.Offset+limit < rootCount,
	}, nil
}
