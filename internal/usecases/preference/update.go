package preference

import (
	"context"
	"encoding/json"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type UpdateUsecase interface {
	Execute(context.Context, string, string, UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
}

type UpdateInput struct{ Settings json.RawMessage }
type UpdateOutput struct {
	Data PreferenceData `json:"data"`
}
type updateUsecase struct{ contextFactory appcontext.Factory }

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, userID, scope string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	if !validScope(scope) {
		return nil, apperrors.NewBadRequestError("invalid preference scope")
	}
	var settings map[string]any
	if len(input.Settings) == 0 || len(input.Settings) > 16*1024 || json.Unmarshal(input.Settings, &settings) != nil || settings == nil {
		return nil, apperrors.NewBadRequestError("settings must be a JSON object")
	}
	canonical, err := json.Marshal(settings)
	if err != nil {
		return nil, apperrors.NewBadRequestError("invalid settings")
	}
	if err := u.contextFactory().Repositories.Preference.Upsert(ctx, userID, scope, canonical); err != nil {
		return nil, apperrors.NewApplicationError(mappings.PreferenceUpdateError, err)
	}
	return &UpdateOutput{Data: PreferenceData{Scope: scope, Settings: canonical}}, nil
}

func validScope(scope string) bool {
	// Add a scope here deliberately when a new UI needs persisted settings.
	return scope == "teacher.dashboard" || scope == "student.dashboard"
}
