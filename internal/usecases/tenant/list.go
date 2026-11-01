package tenant

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type ListUsecase interface {
	Execute(context.Context, ListInput) (*ListOutput, apperrors.ApplicationError)
}

type listUsecase struct {
	factory appcontext.Factory
}

type ListInput struct {
	BearerToken string
}

func NewListUsecase(factory appcontext.Factory) ListUsecase {
	return &listUsecase{factory: factory}
}

func (u *listUsecase) Execute(ctx context.Context, input ListInput) (*ListOutput, apperrors.ApplicationError) {
	app := u.factory()

	tenants, err := app.Repositories.Tenant.List(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TenantListError, err)
	}

	schools, err := app.Integrations.PractiqAPI.ListAllSchools(ctx, input.BearerToken)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TenantScopeUnavailableError, err)
	}
	index := schoolsByID(schools)

	result := make([]TenantAdminData, 0, len(tenants))
	for _, value := range tenants {
		school := index[value.SchoolID]
		result = append(result, TenantAdminData{
			ID:       value.ID,
			SchoolID: value.SchoolID,
			Name:     school.Name,
			Status:   value.Status,
			Eligible: eligible(school),
		})
	}
	return &ListOutput{Data: result}, nil
}
