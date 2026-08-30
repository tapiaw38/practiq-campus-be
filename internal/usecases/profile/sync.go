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
		Email       string
		// BearerToken is the caller's own "Bearer <jwt>" header, forwarded to
		// practiq-be so a practiq "student" profile can override the
		// auth-api-be-role-derived guess (e.g. someone with the global
		// "admin" role who is still a plain student in practiq).
		BearerToken string
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

	profileType := input.ProfileType
	if practiqProfile, err := app.Integrations.PractiqAPI.GetProfile(ctx, input.BearerToken, input.ID); err == nil && practiqProfile != nil {
		if practiqProfile.ProfileType == "student" {
			profileType = "student"
		}
	}

	p := domain.Profile{
		ID:          input.ID,
		ProfileType: profileType,
		FullName:    input.FullName,
		Email:       input.Email,
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
