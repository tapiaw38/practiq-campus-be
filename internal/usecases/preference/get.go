package preference

import (
	"context"
	"encoding/json"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type GetUsecase interface {
	Execute(context.Context, string, string) (*GetOutput, apperrors.ApplicationError)
}

type GetOutput struct {
	Data PreferenceData `json:"data"`
}

type PreferenceData struct {
	Scope    string          `json:"scope"`
	Settings json.RawMessage `json:"settings"`
}

type getUsecase struct{ contextFactory appcontext.Factory }

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{contextFactory: contextFactory}
}

func (u *getUsecase) Execute(ctx context.Context, userID, scope string) (*GetOutput, apperrors.ApplicationError) {
	if !validScope(scope) {
		return nil, apperrors.NewBadRequestError("invalid preference scope")
	}
	settings, _, err := u.contextFactory().Repositories.Preference.Get(ctx, userID, scope)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PreferenceGetError, err)
	}
	return &GetOutput{Data: PreferenceData{Scope: scope, Settings: settings}}, nil
}
