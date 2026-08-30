package profile

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
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
		Search  string
		Page    int
		PerPage int
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
// not just students.
func (u *listStudentsUsecase) Execute(ctx context.Context, input ListStudentsInput) (*ListStudentsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PerPage < 1 || input.PerPage > 20 {
		input.PerPage = 20
	}
	input.Search = strings.TrimSpace(input.Search)

	profiles, total, err := app.Repositories.Profile.ListAll(ctx, input.Search, input.PerPage, (input.Page-1)*input.PerPage)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}

	totalPages := (total + input.PerPage - 1) / input.PerPage
	return &ListStudentsOutput{Data: toProfileDataList(profiles), Meta: ListMeta{Page: input.Page, PerPage: input.PerPage, Total: total, TotalPages: totalPages}}, nil
}
