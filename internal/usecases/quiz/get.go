package quiz

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/unlock"
	"github.com/tapiaw38/practiq-campus-be/internal/usecases/campusaccess"
)

type (
	GetUsecase interface {
		Execute(context.Context, string, bool, string) (*GetOutput, apperrors.ApplicationError)
	}

	getUsecase struct {
		contextFactory appcontext.Factory
	}

	GetOutput struct {
		Data QuizData `json:"data"`
	}
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{contextFactory: contextFactory}
}

// Execute returns one quiz, subject to the same rules the listing applies.
//
// It used to take an id and nothing else, so a direct GET answered for anyone
// who knew the id: it skipped the enrolment check, skipped the group
// restriction that hides a quiz from students outside it, and reported
// locked=false regardless of an unmet prerequisite. Fetching one item is the
// same question as fetching many, and has to survive the same filters.
func (u *getUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	q, err := app.Repositories.Quiz.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.QuizGetError, err)
	}
	if q == nil {
		return nil, apperrors.NewApplicationError(mappings.QuizNotFoundError, nil)
	}

	reader, appErr := campusaccess.EnsureCanReadCourse(ctx, app, requesterID, isSuperAdmin, q.CourseID)
	if appErr != nil {
		return nil, appErr
	}

	status := unlock.Status{}
	if !reader.IsManager {
		if q.VisibleGroupID != nil {
			member, err := app.Repositories.CourseGroup.IsMember(ctx, *q.VisibleGroupID, requesterID)
			if err != nil {
				return nil, apperrors.NewInternalError(err)
			}
			// Restricted to a group this student isn't in: for them it does
			// not exist, the same answer the listing gives by omitting it.
			if !member {
				return nil, apperrors.NewApplicationError(mappings.QuizNotFoundError, nil)
			}
		}
		status, err = unlock.Check(ctx, app.Repositories, q.UnlockAfterType, q.UnlockAfterID, requesterID)
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
	}

	return &GetOutput{Data: toQuizData(*q, status)}, nil
}
