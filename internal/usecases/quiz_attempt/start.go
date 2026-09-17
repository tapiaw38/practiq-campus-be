package quiz_attempt

import (
	"context"
	"errors"
	"time"

	attemptrepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/quiz_attempt"
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
		// Resumed tells the client it is rejoining an attempt already under
		// way rather than starting a fresh one, so it can say so instead of
		// implying the clock just began.
		Resumed bool `json:"resumed"`
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

	questions, err := app.Repositories.Quiz.ListQuestions(ctx, quizID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if len(questions) == 0 {
		return nil, apperrors.NewBadRequestError("this quiz has no questions yet")
	}

	// One atomic step: it hands back the attempt already in progress if there
	// is one, so reloading the page rejoins the exam instead of spending
	// another try on it, and it settles the attempt number under a lock so
	// two clicks cannot collide on it.
	result, err := app.Repositories.QuizAttempt.StartOrResume(ctx, attemptrepo.StartParams{
		QuizID:        quizID,
		UserID:        requesterID,
		MaxAttempts:   q.MaxAttempts,
		TimeLimitSecs: q.TimeLimitSecs,
		HardDeadline:  q.AvailableUntil,
	})
	if errors.Is(err, attemptrepo.ErrAttemptsExhausted) {
		return nil, apperrors.NewApplicationError(mappings.QuizAttemptsExhaustedError, nil)
	}
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}

	return &StartOutput{
		Attempt:   toAttemptData(result.Attempt),
		Questions: toStudentQuestionData(questions),
		TimeLimit: q.TimeLimitSecs,
		Resumed:   result.Resumed,
	}, nil
}
