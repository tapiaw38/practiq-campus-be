package tenant

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type ListMineUsecase interface {
	Execute(context.Context, ListMineInput) (*ListMineOutput, apperrors.ApplicationError)
}

type listMineUsecase struct {
	factory appcontext.Factory
}

type ListMineInput struct {
	BearerToken  string
	IsSuperAdmin bool
}

func NewListMineUsecase(factory appcontext.Factory) ListMineUsecase {
	return &listMineUsecase{factory: factory}
}

func (u *listMineUsecase) Execute(ctx context.Context, input ListMineInput) (*ListMineOutput, apperrors.ApplicationError) {
	app := u.factory()

	if input.IsSuperAdmin {
		return u.listAllForSuperAdmin(ctx, input)
	}

	schools, err := app.Integrations.PractiqAPI.ListMySchools(ctx, input.BearerToken)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TenantMembershipUnavailableError, err)
	}

	result := make([]TenantData, 0, len(schools))
	for _, school := range schools {
		if !eligible(school) {
			continue
		}
		value, err := app.Repositories.Tenant.GetBySchoolID(ctx, school.ID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.TenantListError, err)
		}
		if value != nil && value.Status == domain.TenantStatusActive {
			result = append(result, TenantData{ID: value.ID, SchoolID: school.ID, Name: school.Name, Role: school.Role})
		}
	}
	return &ListMineOutput{Data: result}, nil
}

func (u *listMineUsecase) listAllForSuperAdmin(ctx context.Context, input ListMineInput) (*ListMineOutput, apperrors.ApplicationError) {
	app := u.factory()

	tenants, err := app.Repositories.Tenant.List(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TenantListError, err)
	}

	names := map[string]string{}
	if schools, err := app.Integrations.PractiqAPI.ListAllSchools(ctx, input.BearerToken); err == nil {
		for _, school := range schools {
			names[school.ID] = school.Name
		}
	}

	result := make([]TenantData, 0, len(tenants))
	for _, value := range tenants {
		if value.Status != domain.TenantStatusActive {
			continue
		}
		result = append(result, TenantData{ID: value.ID, SchoolID: value.SchoolID, Name: names[value.SchoolID]})
	}
	return &ListMineOutput{Data: result}, nil
}
