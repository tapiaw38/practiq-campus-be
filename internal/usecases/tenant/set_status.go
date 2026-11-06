package tenant

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type SetStatusUsecase interface {
	Execute(context.Context, SetStatusInput) (*SetStatusOutput, apperrors.ApplicationError)
}

type setStatusUsecase struct {
	factory appcontext.Factory
}

type SetStatusInput struct {
	BearerToken string
	TenantID    string
	Status      string
}

func NewSetStatusUsecase(factory appcontext.Factory) SetStatusUsecase {
	return &setStatusUsecase{factory: factory}
}

func (u *setStatusUsecase) Execute(ctx context.Context, input SetStatusInput) (*SetStatusOutput, apperrors.ApplicationError) {
	app := u.factory()

	existing, err := app.Repositories.Tenant.Get(ctx, input.TenantID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TenantLookupError, err)
	}
	if existing == nil {
		return nil, apperrors.NewApplicationError(mappings.TenantNotFoundError, nil)
	}

	schools, schoolsErr := app.Integrations.PractiqAPI.ListAllSchools(ctx, input.BearerToken)
	if input.Status == domain.TenantStatusActive {
		if schoolsErr != nil {
			return nil, apperrors.NewApplicationError(mappings.TenantScopeUnavailableError, schoolsErr)
		}
		school, found := schoolsByID(schools)[existing.SchoolID]
		if !found || !eligible(school) {
			return nil, apperrors.NewApplicationError(mappings.TenantSchoolNoLongerEligibleError, nil)
		}
	}

	value, err := app.Repositories.Tenant.SetStatus(ctx, existing.ID, input.Status)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TenantStatusUpdateError, err)
	}

	data := TenantAdminData{ID: value.ID, SchoolID: value.SchoolID, Status: value.Status}
	if schoolsErr == nil {
		if school, found := schoolsByID(schools)[existing.SchoolID]; found {
			name := school.Name
			qualifies := eligible(school)
			data.Name = &name
			data.Eligible = &qualifies
		}
	}
	return &SetStatusOutput{Data: data}, nil
}
