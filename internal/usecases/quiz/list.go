package quiz

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(context.Context, string) (*ListOutput, apperrors.ApplicationError)
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

func (u *listUsecase) Execute(ctx context.Context, courseID string) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	quizzes, err := app.Repositories.Quiz.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.QuizListError, err)
	}

	return &ListOutput{Data: toQuizDataList(quizzes)}, nil
}
