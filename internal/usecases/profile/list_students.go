package profile

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/identity"
)

type (
	ListStudentsUsecase interface {
		Execute(context.Context, ListStudentsInput) (*ListStudentsOutput, apperrors.ApplicationError)
	}

	listStudentsUsecase struct {
		contextFactory appcontext.Factory
	}

	ListStudentsOutput struct {
		Data []ProfileData `json:"data"`
		Meta ListMeta      `json:"meta"`
	}
	ListStudentsInput struct {
		Search      string
		Page        int
		PerPage     int
		BearerToken string
	}
	ListMeta struct {
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	}
)

func NewListStudentsUsecase(contextFactory appcontext.Factory) ListStudentsUsecase {
	return &listStudentsUsecase{contextFactory: contextFactory}
}

// Execute lists every locally known Campus profile — both students and
// teachers — since the admin "Usuarios" screen manages accounts in general,
// not just students. Name/email no longer live locally, so search and
// pagination both happen in memory after a single batch identity lookup.
// ponytail: fine at classroom-app scale (dozens/hundreds of accounts); if
// the user base grows large, move search to an auth-api-be endpoint that
// filters server-side instead of fetching every local profile per request.
func (u *listStudentsUsecase) Execute(ctx context.Context, input ListStudentsInput) (*ListStudentsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PerPage < 1 || input.PerPage > 20 {
		input.PerPage = 20
	}
	search := strings.ToLower(strings.TrimSpace(input.Search))

	profiles, _, err := app.Repositories.Profile.ListAll(ctx, 1_000_000, 0)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}

	ids := make([]string, 0, len(profiles))
	for _, p := range profiles {
		ids = append(ids, p.ID)
	}
	names, err := identity.Names(ctx, app.Integrations.AuthAPI, input.BearerToken, ids)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}

	matches := make([]ProfileData, 0, len(profiles))
	for _, p := range profiles {
		info := names[p.ID]
		fullName := identity.FullName(info, p.ID)
		if search != "" &&
			!strings.Contains(strings.ToLower(fullName), search) &&
			!strings.Contains(strings.ToLower(info.Email), search) {
			continue
		}
		matches = append(matches, toProfileData(p, fullName, info.Email))
	}

	total := len(matches)
	totalPages := (total + input.PerPage - 1) / input.PerPage
	start := (input.Page - 1) * input.PerPage
	if start > total {
		start = total
	}
	end := start + input.PerPage
	if end > total {
		end = total
	}

	return &ListStudentsOutput{Data: matches[start:end], Meta: ListMeta{Page: input.Page, PerPage: input.PerPage, Total: total, TotalPages: totalPages}}, nil
}
