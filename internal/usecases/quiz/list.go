package quiz

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/unlock"
)

type (
	ListUsecase interface {
		Execute(context.Context, string, bool, string) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []QuizData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

// Execute filters out quizzes restricted to a group the requester isn't in,
// and flags the rest as locked when their prerequisite isn't met yet — but
// only for a plain student. The course owner and a superadmin manage every
// item, so nothing is ever hidden or locked from them.
func (u *listUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	quizzes, err := app.Repositories.Quiz.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.QuizListError, err)
	}

	isManager := isSuperAdmin
	if !isManager {
		course, err := app.Repositories.Course.Get(ctx, courseID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
		}
		isManager = course != nil && course.OwnerID == requesterID
	}

	out := make([]QuizData, 0, len(quizzes))
	for _, q := range quizzes {
		if !isManager && q.VisibleGroupID != nil {
			member, err := app.Repositories.CourseGroup.IsMember(ctx, *q.VisibleGroupID, requesterID)
			if err != nil {
				return nil, apperrors.NewInternalError(err)
			}
			if !member {
				continue
			}
		}
		status := unlock.Status{}
		if !isManager {
			status, err = unlock.Check(ctx, app.Repositories, q.UnlockAfterType, q.UnlockAfterID, requesterID)
			if err != nil {
				return nil, apperrors.NewInternalError(err)
			}
		}
		out = append(out, toQuizData(q, status))
	}

	return &ListOutput{Data: out}, nil
}
