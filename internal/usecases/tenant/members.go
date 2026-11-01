package tenant

import (
	"context"
	"errors"
	"net/http"

	"github.com/tapiaw38/practiq-campus-be/internal/adapters/web/integrations/practiqapi"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

func mapUpstreamError(err error, fallbackCode, fallbackMessage string) apperrors.ApplicationError {
	var upstream *practiqapi.UpstreamError
	if errors.As(err, &upstream) && upstream.Status >= 400 && upstream.Status < 500 {
		return apperrors.NewApplicationError(mappings.ErrorDetails{
			InternalCode: upstream.Code,
			StatusCode:   upstream.Status,
			Message:      upstream.Message,
		}, err)
	}
	return apperrors.NewApplicationError(mappings.ErrorDetails{
		InternalCode: fallbackCode,
		StatusCode:   http.StatusBadGateway,
		Message:      fallbackMessage,
	}, err)
}

type (
	ListMembersUsecase interface {
		Execute(context.Context, ListMembersInput) (*ListMembersOutput, apperrors.ApplicationError)
	}
	AddMemberUsecase interface {
		Execute(context.Context, AddMemberInput) apperrors.ApplicationError
	}
	RemoveMemberUsecase interface {
		Execute(context.Context, RemoveMemberInput) apperrors.ApplicationError
	}

	listMembersUsecase  struct{ factory appcontext.Factory }
	addMemberUsecase    struct{ factory appcontext.Factory }
	removeMemberUsecase struct{ factory appcontext.Factory }

	ListMembersInput struct {
		BearerToken string
		TenantID    string
	}
	AddMemberInput struct {
		BearerToken string
		TenantID    string
		UserID      string
		Role        string
	}
	RemoveMemberInput struct {
		BearerToken string
		TenantID    string
		UserID      string
	}
)

func NewListMembersUsecase(factory appcontext.Factory) ListMembersUsecase {
	return &listMembersUsecase{factory: factory}
}
func NewAddMemberUsecase(factory appcontext.Factory) AddMemberUsecase {
	return &addMemberUsecase{factory: factory}
}
func NewRemoveMemberUsecase(factory appcontext.Factory) RemoveMemberUsecase {
	return &removeMemberUsecase{factory: factory}
}

func (u *listMembersUsecase) Execute(ctx context.Context, input ListMembersInput) (*ListMembersOutput, apperrors.ApplicationError) {
	app := u.factory()

	value, err := app.Repositories.Tenant.Get(ctx, input.TenantID)
	if err != nil || value == nil {
		return nil, apperrors.NewApplicationError(mappings.TenantMemberInstitutionNotFoundError, err)
	}

	members, err := app.Integrations.PractiqAPI.ListSchoolMembers(ctx, input.BearerToken, value.SchoolID)
	if err != nil {
		return nil, mapUpstreamError(err, "tenant:members-unavailable", "could not load institution members")
	}
	return &ListMembersOutput{Data: members}, nil
}

func (u *addMemberUsecase) Execute(ctx context.Context, input AddMemberInput) apperrors.ApplicationError {
	app := u.factory()

	value, err := app.Repositories.Tenant.Get(ctx, input.TenantID)
	if err != nil || value == nil {
		return apperrors.NewApplicationError(mappings.TenantMemberInstitutionNotFoundError, err)
	}

	if err := app.Integrations.PractiqAPI.AddSchoolMember(ctx, input.BearerToken, value.SchoolID, input.UserID, input.Role); err != nil {
		return mapUpstreamError(err, "tenant:add-member-error", "could not add member to institution")
	}
	return nil
}

func (u *removeMemberUsecase) Execute(ctx context.Context, input RemoveMemberInput) apperrors.ApplicationError {
	app := u.factory()

	value, err := app.Repositories.Tenant.Get(ctx, input.TenantID)
	if err != nil || value == nil {
		return apperrors.NewApplicationError(mappings.TenantMemberInstitutionNotFoundError, err)
	}

	if err := app.Integrations.PractiqAPI.RemoveSchoolMember(ctx, input.BearerToken, value.SchoolID, input.UserID); err != nil {
		return mapUpstreamError(err, "tenant:remove-member-error", "could not remove member from institution")
	}
	return nil
}
