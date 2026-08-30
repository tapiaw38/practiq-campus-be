package profile

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	ListMyPractiqStudentsUsecase interface {
		Execute(context.Context, string) (*ListMyPractiqStudentsOutput, apperrors.ApplicationError)
	}

	listMyPractiqStudentsUsecase struct {
		contextFactory appcontext.Factory
	}

	PractiqStudentData struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	ListMyPractiqStudentsOutput struct {
		Data []PractiqStudentData `json:"data"`
	}
)

func NewListMyPractiqStudentsUsecase(contextFactory appcontext.Factory) ListMyPractiqStudentsUsecase {
	return &listMyPractiqStudentsUsecase{contextFactory: contextFactory}
}

// Execute is a pure pass-through to practiq-be, scoped to the caller's own
// bearer token — it can only ever return that same teacher's own students,
// same as calling practiq-be directly would.
func (u *listMyPractiqStudentsUsecase) Execute(ctx context.Context, bearerToken string) (*ListMyPractiqStudentsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	students, err := app.Integrations.PractiqAPI.ListMyStudents(ctx, bearerToken)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PractiqLookupError, err)
	}

	data := make([]PractiqStudentData, 0, len(students))
	for _, s := range students {
		data = append(data, PractiqStudentData{ID: s.ID, Name: s.Name, Email: s.Email})
	}

	return &ListMyPractiqStudentsOutput{Data: data}, nil
}
