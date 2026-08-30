package submission

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	GetMineUsecase interface {
		Execute(context.Context, string, string) (*GetMineOutput, apperrors.ApplicationError)
	}

	getMineUsecase struct {
		contextFactory appcontext.Factory
	}

	GetMineOutput struct {
		Data *SubmissionData `json:"data"`
	}
)

func NewGetMineUsecase(contextFactory appcontext.Factory) GetMineUsecase {
	return &getMineUsecase{contextFactory: contextFactory}
}

// Execute returns nil Data (not an error) when the requester hasn't
// submitted yet — "no submission" is a normal state for an assignment, not
// a failure.
func (u *getMineUsecase) Execute(ctx context.Context, requesterID, assignmentID string) (*GetMineOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	s, err := app.Repositories.Submission.GetByAssignmentAndUser(ctx, assignmentID, requesterID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubmissionGetError, err)
	}
	if s == nil {
		return &GetMineOutput{Data: nil}, nil
	}

	data := toSubmissionData(*s)
	return &GetMineOutput{Data: &data}, nil
}
