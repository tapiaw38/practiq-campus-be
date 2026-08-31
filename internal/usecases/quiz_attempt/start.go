package quiz_attempt

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/unlock"
)

type (
	StartUsecase interface {
		Execute(ctx context.Context, requesterID, quizID string) (*StartOutput, apperrors.ApplicationError)
	}

	startUsecase struct {
		contextFactory appcontext.Factory
	}

	StartOutput struct {
		Attempt   AttemptData           `json:"attempt"`
		Questions []StudentQuestionData `json:"questions"`
		TimeLimit *int                  `json:"time_limit_secs"`
	}
)

func NewStartUsecase(contextFactory appcontext.Factory) StartUsecase {
	return &startUsecase{contextFactory: contextFactory}
}

func (u *startUsecase) Execute(ctx context.Context, requesterID, quizID string) (*StartOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	q, err := app.Repositories.Quiz.Get(ctx, quizID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.QuizGetError, err)
	}
	if q == nil {
		return nil, apperrors.NewApplicationError(mappings.QuizNotFoundError, nil)
	}

	enrollment, err := app.Repositories.Enrollment.GetByCourseAndUser(ctx, q.CourseID, requesterID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if enrollment == nil {
		return nil, apperrors.NewApplicationError(mappings.QuizNotEnrolledError, nil)
	}
	if q.VisibleGroupID != nil {
		member, err := app.Repositories.CourseGroup.IsMember(ctx, *q.VisibleGroupID, requesterID)
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		if !member {
			return nil, apperrors.NewForbiddenError()
		}
	}
	if status, err := unlock.Check(ctx, app.Repositories, q.UnlockAfterType, q.UnlockAfterID, requesterID); err != nil {
		return nil, apperrors.NewInternalError(err)
	} else if status.Locked {
		return nil, apperrors.NewBadRequestError(status.Reason)
	}

	now := time.Now()
	if (q.ScheduledAt != nil && now.Before(*q.ScheduledAt)) || (q.AvailableUntil != nil && now.After(*q.AvailableUntil)) {
		return nil, apperrors.NewApplicationError(mappings.QuizNotAvailableError, nil)
	}

	count, err := app.Repositories.QuizAttempt.CountByUser(ctx, quizID, requesterID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if q.MaxAttempts > 0 && count >= q.MaxAttempts {
		return nil, apperrors.NewApplicationError(mappings.QuizAttemptsExhaustedError, nil)
	}

	questions, err := app.Repositories.Quiz.ListQuestions(ctx, quizID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if len(questions) == 0 {
		return nil, apperrors.NewBadRequestError("this quiz has no questions yet")
	}

	id, err := app.Repositories.QuizAttempt.Create(ctx, domain.QuizAttempt{
		QuizID:        quizID,
		UserID:        requesterID,
		AttemptNumber: count + 1,
	})
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}

	created, err := app.Repositories.QuizAttempt.Get(ctx, id)
	if err != nil || created == nil {
		return nil, apperrors.NewInternalError(err)
	}

	return &StartOutput{
		Attempt:   toAttemptData(*created),
		Questions: toStudentQuestionData(questions),
		TimeLimit: q.TimeLimitSecs,
	}, nil
}
