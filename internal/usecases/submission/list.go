package submission

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	ListByAssignmentUsecase interface {
		Execute(context.Context, string, bool, string) (*ListOutput, apperrors.ApplicationError)
	}

	listByAssignmentUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []SubmissionData `json:"data"`
	}
)

func NewListByAssignmentUsecase(contextFactory appcontext.Factory) ListByAssignmentUsecase {
	return &listByAssignmentUsecase{contextFactory: contextFactory}
}

func (u *listByAssignmentUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, assignmentID string) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	a, err := app.Repositories.Assignment.Get(ctx, assignmentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentGetError, err)
	}
	if a == nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentNotFoundError, nil)
	}

	if !isSuperAdmin {
		c, err := app.Repositories.Course.Get(ctx, a.CourseID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
		}
		if c == nil || c.OwnerID != requesterID {
			return nil, apperrors.NewForbiddenError()
		}
	}

	submissions, err := app.Repositories.Submission.ListByAssignment(ctx, assignmentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubmissionListError, err)
	}

	return &ListOutput{Data: toSubmissionDataList(submissions)}, nil
}
