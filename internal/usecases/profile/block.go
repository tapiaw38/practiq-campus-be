package profile

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/identity"
)

type BlockUsecase interface {
	Execute(ctx context.Context, requesterID, userID string, blocked bool, bearerToken string) (*BlockOutput, apperrors.ApplicationError)
}
type BlockOutput struct {
	Data ProfileData `json:"data"`
}
type blockUsecase struct{ contextFactory appcontext.Factory }

func NewBlockUsecase(contextFactory appcontext.Factory) BlockUsecase {
	return &blockUsecase{contextFactory: contextFactory}
}
func (u *blockUsecase) Execute(ctx context.Context, requesterID, userID string, blocked bool, bearerToken string) (*BlockOutput, apperrors.ApplicationError) {
	if requesterID == userID {
		return nil, apperrors.NewBadRequestError("cannot block yourself")
	}
	app := u.contextFactory()
	profile, err := app.Repositories.Profile.Get(ctx, userID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if profile == nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileNotFoundError, nil)
	}
	if err := app.Repositories.Profile.SetBlocked(ctx, userID, blocked); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileSyncError, err)
	}
	profile.IsBlocked = blocked

	names, err := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, []string{userID})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	info := names[userID]

	return &BlockOutput{Data: toProfileData(*profile, identity.FullName(info, userID), info.Email)}, nil
}
