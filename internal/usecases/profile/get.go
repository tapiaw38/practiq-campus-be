package profile

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	GetUsecase interface {
		Execute(context.Context, string) (*GetOutput, apperrors.ApplicationError)
	}

	getUsecase struct {
		contextFactory appcontext.Factory
	}

	GetOutput struct {
		Data ProfileData `json:"data"`
	}
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{contextFactory: contextFactory}
}

func (u *getUsecase) Execute(ctx context.Context, userID string) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	p, err := app.Repositories.Profile.Get(ctx, userID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if p == nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileNotFoundError, nil)
	}

	return &GetOutput{Data: toProfileData(*p)}, nil
}
