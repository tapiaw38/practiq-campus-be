package tenant

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type ActivateUsecase interface {
	Execute(context.Context, ActivateInput) (*ActivateOutput, apperrors.ApplicationError)
}

type activateUsecase struct {
	factory appcontext.Factory
}

type ActivateInput struct {
	BearerToken string
	SchoolID    string
}

func NewActivateUsecase(factory appcontext.Factory) ActivateUsecase {
	return &activateUsecase{factory: factory}
}

func (u *activateUsecase) Execute(ctx context.Context, input ActivateInput) (*ActivateOutput, apperrors.ApplicationError) {
	app := u.factory()

	schools, err := app.Integrations.PractiqAPI.ListAllSchools(ctx, input.BearerToken)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TenantScopeUnavailableError, err)
	}
	school, found := schoolsByID(schools)[input.SchoolID]
	if !found {
		return nil, apperrors.NewApplicationError(mappings.TenantSchoolNotFoundError, nil)
	}
	if !eligible(school) {
		return nil, apperrors.NewApplicationError(mappings.TenantSchoolNotEligibleError, nil)
	}

	value, err := app.Repositories.Tenant.Create(ctx, domain.Tenant{SchoolID: school.ID, Status: domain.TenantStatusActive})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TenantActivateError, err)
	}
	return &ActivateOutput{Data: TenantAdminData{ID: value.ID, SchoolID: value.SchoolID, Name: school.Name, Status: value.Status, Eligible: true}}, nil
}
