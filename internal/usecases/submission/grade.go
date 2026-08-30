package submission

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	GradeUsecase interface {
		Execute(context.Context, string, bool, string, GradeInput) (*GradeOutput, apperrors.ApplicationError)
	}

	gradeUsecase struct {
		contextFactory appcontext.Factory
	}

	GradeInput struct {
		Score    int
		Feedback string
	}

	GradeOutput struct {
		Data SubmissionData `json:"data"`
	}
)

func NewGradeUsecase(contextFactory appcontext.Factory) GradeUsecase {
	return &gradeUsecase{contextFactory: contextFactory}
}

func (u *gradeUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, submissionID string, input GradeInput) (*GradeOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	s, err := app.Repositories.Submission.Get(ctx, submissionID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubmissionGetError, err)
	}
	if s == nil {
		return nil, apperrors.NewApplicationError(mappings.SubmissionNotFoundError, nil)
	}

	a, err := app.Repositories.Assignment.Get(ctx, s.AssignmentID)
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

	if input.Score < 0 || input.Score > a.MaxScore {
		return nil, apperrors.NewBadRequestError("score must be between 0 and the assignment's max score")
	}

	if err := app.Repositories.Submission.Grade(ctx, submissionID, input.Score, input.Feedback); err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubmissionCreateError, err)
	}

	updated, err := app.Repositories.Submission.Get(ctx, submissionID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubmissionGetError, err)
	}
	if updated == nil {
		return nil, apperrors.NewInternalError(nil)
	}

	return &GradeOutput{Data: toSubmissionData(*updated)}, nil
}
