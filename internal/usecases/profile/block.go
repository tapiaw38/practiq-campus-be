package profile

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type BlockUsecase interface {
	Execute(context.Context, string, string, bool) (*BlockOutput, apperrors.ApplicationError)
}
type BlockOutput struct {
	Data ProfileData `json:"data"`
}
type blockUsecase struct{ contextFactory appcontext.Factory }

func NewBlockUsecase(contextFactory appcontext.Factory) BlockUsecase {
	return &blockUsecase{contextFactory: contextFactory}
}
func (u *blockUsecase) Execute(ctx context.Context, requesterID, userID string, blocked bool) (*BlockOutput, apperrors.ApplicationError) {
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
	return &BlockOutput{Data: toProfileData(*profile)}, nil
}
