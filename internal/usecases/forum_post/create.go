package forum_post

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/identity"
)

type (
	CreateUsecase interface {
		Execute(context.Context, string, bool, string, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		Body        string
		ParentID    *string
		BearerToken string
	}

	CreateOutput struct {
		Data PostData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, threadID string, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	if strings.TrimSpace(input.Body) == "" {
		return nil, apperrors.NewBadRequestError("body is required")
	}

	app := u.contextFactory()

	if _, appErr := requesterCanAccessThread(ctx, app, requesterID, isSuperAdmin, threadID); appErr != nil {
		return nil, appErr
	}
	if input.ParentID != nil {
		parent, err := app.Repositories.ForumPost.Get(ctx, *input.ParentID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.PostListError, err)
		}
		if parent == nil || parent.ThreadID != threadID {
			return nil, apperrors.NewBadRequestError("parent post does not belong to thread")
		}
	}

	id, err := app.Repositories.ForumPost.Create(ctx, domain.ForumPost{
		ThreadID: threadID,
		ParentID: input.ParentID,
		AuthorID: requesterID,
		Body:     input.Body,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PostCreateError, err)
	}

	created, err := app.Repositories.ForumPost.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PostListError, err)
	}
	if created == nil {
		return nil, apperrors.NewInternalError(nil)
	}

	names, err := identity.Names(ctx, app.Integrations.AuthAPI, input.BearerToken, []string{created.AuthorID})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PostListError, err)
	}
	created.AuthorName = identity.FullName(names[created.AuthorID], created.AuthorID)

	return &CreateOutput{Data: toPostData(*created)}, nil
}
