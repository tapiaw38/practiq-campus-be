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
