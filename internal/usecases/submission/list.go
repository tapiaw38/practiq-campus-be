package submission

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/identity"
)

type (
	ListByAssignmentUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, assignmentID, bearerToken string) (*ListOutput, apperrors.ApplicationError)
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

func (u *listByAssignmentUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, assignmentID, bearerToken string) (*ListOutput, apperrors.ApplicationError) {
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

	ids := make([]string, 0, len(submissions))
	for _, s := range submissions {
		ids = append(ids, s.UserID)
	}
	names, err := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, ids)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	for i, s := range submissions {
		submissions[i].UserName = identity.FullName(names[s.UserID], s.UserID)
	}

	data := toSubmissionDataList(submissions)
	for i, submission := range submissions {
		scores, err := app.Repositories.Rubric.Scores(ctx, submission.ID)
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		data[i] = withAttachments(app, withRubricScores(data[i], scores))
	}
	return &ListOutput{Data: data}, nil
}
