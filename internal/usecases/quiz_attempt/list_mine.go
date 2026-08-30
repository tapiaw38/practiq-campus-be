package quiz_attempt

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
)

type (
	ListMineUsecase interface {
		Execute(ctx context.Context, requesterID, quizID string) (*ListOutput, apperrors.ApplicationError)
	}

	listMineUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []AttemptData `json:"data"`
	}
)

func NewListMineUsecase(contextFactory appcontext.Factory) ListMineUsecase {
	return &listMineUsecase{contextFactory: contextFactory}
}

func (u *listMineUsecase) Execute(ctx context.Context, requesterID, quizID string) (*ListOutput, apperrors.ApplicationError) {
	attempts, err := u.contextFactory().Repositories.QuizAttempt.ListMine(ctx, quizID, requesterID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return &ListOutput{Data: toAttemptDataList(attempts)}, nil
}
