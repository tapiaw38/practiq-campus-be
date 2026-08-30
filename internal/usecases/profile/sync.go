package profile

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	SyncUsecase interface {
		Execute(context.Context, SyncInput) (*SyncOutput, apperrors.ApplicationError)
	}

	syncUsecase struct {
		contextFactory appcontext.Factory
	}

	SyncInput struct {
		ID          string
		ProfileType string
		FullName    string
	}

	SyncOutput struct {
		Data ProfileData `json:"data"`
	}
)

func NewSyncUsecase(contextFactory appcontext.Factory) SyncUsecase {
	return &syncUsecase{contextFactory: contextFactory}
}

func (u *syncUsecase) Execute(ctx context.Context, input SyncInput) (*SyncOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	p := domain.Profile{
		ID:          input.ID,
		ProfileType: input.ProfileType,
		FullName:    input.FullName,
	}

	if err := app.Repositories.Profile.Upsert(ctx, p); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileSyncError, err)
	}

	updated, err := app.Repositories.Profile.Get(ctx, input.ID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if updated == nil {
		return nil, apperrors.NewInternalError(nil)
	}

	return &SyncOutput{Data: toProfileData(*updated)}, nil
}
