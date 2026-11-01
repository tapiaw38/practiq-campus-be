package profile

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/authapi"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	CreateOrSyncUserUsecase interface {
		Execute(context.Context, CreateOrSyncUserInput) (*SyncOutput, apperrors.ApplicationError)
	}

	createOrSyncUserUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateOrSyncUserInput struct {
		BearerToken string
		Email       string
		FirstName   string
		LastName    string
		Password    string
	}
)

func NewProvisionStudentUsecase(contextFactory appcontext.Factory) CreateOrSyncUserUsecase {
	return &createOrSyncUserUsecase{contextFactory: contextFactory}
}

func (u *createOrSyncUserUsecase) Execute(ctx context.Context, input CreateOrSyncUserInput) (*SyncOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	authUser, err := app.Integrations.AuthAPI.GetByEmail(ctx, input.BearerToken, input.Email)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileProvisionAuthError, err)
	}

	var id, fullName, profileType string

	if authUser != nil {

		id = authUser.Username

		if existing, err := app.Repositories.Profile.Get(ctx, id); err != nil {
			return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
		} else if existing != nil {
			return &SyncOutput{Data: toProfileData(*existing, strings.TrimSpace(authUser.FirstName+" "+authUser.LastName), authUser.Email)}, nil
		}

		fullName = authUser.FirstName + " " + authUser.LastName
		profileType = "student"
		for _, role := range authUser.Roles {
			if role == "admin" || role == "superadmin" {
				profileType = "teacher"
			}
		}

		if practiqProfile, err := app.Integrations.PractiqAPI.GetProfile(ctx, input.BearerToken, id); err == nil && practiqProfile != nil {
			if practiqProfile.ProfileType == "student" {
				profileType = "student"
			}
		}
	} else {
		if input.FirstName == "" || input.LastName == "" || input.Password == "" {
			return nil, apperrors.NewApplicationError(mappings.ProfileMissingFieldsError, nil)
		}

		registered, err := app.Integrations.AuthAPI.Register(ctx, authapi.RegisterInput{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Password:  input.Password,
		})
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.ProfileProvisionAuthError, err)
		}

		id = registered.Username
		fullName = input.FirstName + " " + input.LastName
		profileType = "student"
	}

	p := domain.Profile{
		ID:          id,
		ProfileType: profileType,
	}
	if err := app.Repositories.Profile.Upsert(ctx, p); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileSyncError, err)
	}

	created, err := app.Repositories.Profile.Get(ctx, p.ID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if created == nil {
		return nil, apperrors.NewInternalError(nil)
	}

	return &SyncOutput{Data: toProfileData(*created, fullName, input.Email)}, nil
}
