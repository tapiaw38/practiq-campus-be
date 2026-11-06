package tenant

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type ListEligibleSchoolsUsecase interface {
	Execute(context.Context, ListEligibleSchoolsInput) (*ListEligibleSchoolsOutput, apperrors.ApplicationError)
}

type listEligibleSchoolsUsecase struct {
	factory appcontext.Factory
}

type ListEligibleSchoolsInput struct {
	BearerToken string
}

func NewListEligibleSchoolsUsecase(factory appcontext.Factory) ListEligibleSchoolsUsecase {
	return &listEligibleSchoolsUsecase{factory: factory}
}

func (u *listEligibleSchoolsUsecase) Execute(ctx context.Context, input ListEligibleSchoolsInput) (*ListEligibleSchoolsOutput, apperrors.ApplicationError) {
	app := u.factory()

	schools, err := app.Integrations.PractiqAPI.ListAllSchools(ctx, input.BearerToken)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TenantScopeUnavailableError, err)
	}

	tenants, err := app.Repositories.Tenant.List(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TenantListError, err)
	}

	enabled := make(map[string]bool, len(tenants))
	for _, value := range tenants {
		enabled[value.SchoolID] = true
	}

	result := make([]EligibleSchoolData, 0)
	for _, school := range schools {
		if eligible(school) && !enabled[school.ID] {
			result = append(result, EligibleSchoolData{ID: school.ID, Name: school.Name})
		}
	}
	return &ListEligibleSchoolsOutput{Data: result}, nil
}
