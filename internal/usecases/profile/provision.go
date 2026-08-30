package profile

import (
	"context"

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
		// BearerToken is the requesting superadmin's own "Bearer <jwt>"
		// header, forwarded as-is to auth-api-be/practiq-be — those tokens
		// are valid on every service since they all share one JWT secret,
		// so no separate service-to-service credential is needed.
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

// Execute is superadmin-only. Three paths:
//  1. A campus_profiles row for this email already exists — return it as-is,
//     nothing to do.
//  2. The email already has a shared identity on auth-api-be — reuse it (no
//     new password, no register call). profile_type is derived from that
//     account's auth-api-be roles; when it resolves to "student", it also
//     asks practiq-be for that same user's practiq profile to reuse its
//     display name if practiq-be already has one (student case only, per
//     scope — teacher-in-practiq reconciliation is a later problem).
//  3. Brand-new email — requires first_name/last_name/password and goes
//     through the normal auth-api-be register flow.
func (u *createOrSyncUserUsecase) Execute(ctx context.Context, input CreateOrSyncUserInput) (*SyncOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if existing, err := app.Repositories.Profile.GetByEmail(ctx, input.Email); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	} else if existing != nil {
		return &SyncOutput{Data: toProfileData(*existing)}, nil
	}

	authUser, err := app.Integrations.AuthAPI.GetByEmail(ctx, input.BearerToken, input.Email)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileProvisionAuthError, err)
	}

	var id, fullName, profileType string

	if authUser != nil {
		// campus_profiles.id must match the JWT's user_id claim, which
		// auth-api-be sets to the username (not the UUID id) — self-sync on
		// login keys profiles the same way, so provisioning has to use the
		// same identifier or the two paths create separate rows for the
		// same person.
		id = authUser.Username
		fullName = (authUser.FirstName + " " + authUser.LastName)
		profileType = "student"
		for _, role := range authUser.Roles {
			if role == "admin" || role == "superadmin" {
				profileType = "teacher"
			}
		}

		// practiq-be is the source of truth when it knows this person as a
		// student — a superadmin/teacher account on auth-api-be doesn't mean
		// they aren't also a plain student in practiq. Teacher-in-practiq
		// reconciliation is intentionally out of scope for now, so a
		// practiq "teacher" answer doesn't override anything here.
		if practiqProfile, err := app.Integrations.PractiqAPI.GetProfile(ctx, input.BearerToken, id); err == nil && practiqProfile != nil {
			if practiqProfile.ProfileType == "student" {
				profileType = "student"
				if practiqProfile.Name != "" {
					fullName = practiqProfile.Name
				}
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
		FullName:    fullName,
		Email:       input.Email,
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

	return &SyncOutput{Data: toProfileData(*created)}, nil
}
