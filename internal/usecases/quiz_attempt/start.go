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

		Resumed bool `json:"resumed"`

		SavedAnswers []DraftAnswerData `json:"saved_answers"`

		ServerNow string `json:"server_now"`
	}

	DraftAnswerData struct {
		QuestionID string `json:"question_id"`
		AnswerText string `json:"answer_text"`
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

	windowOpen := q.AcceptsNewAttempts(time.Now())

	questions, err := app.Repositories.Quiz.ListQuestions(ctx, quizID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if len(questions) == 0 {
		return nil, apperrors.NewBadRequestError("this quiz has no questions yet")
	}

	result, err := app.Repositories.QuizAttempt.StartOrResume(ctx, attemptrepo.StartParams{
		QuizID:        quizID,
		UserID:        requesterID,
		MaxAttempts:   q.MaxAttempts,
		TimeLimitSecs: q.TimeLimitSecs,
		HardDeadline:  q.AvailableUntil,
		WindowOpen:    windowOpen,
	})
	switch {
	case errors.Is(err, attemptrepo.ErrAttemptsExhausted):
		return nil, apperrors.NewApplicationError(mappings.QuizAttemptsExhaustedError, nil)
	case errors.Is(err, attemptrepo.ErrQuizWindowClosed):
		return nil, apperrors.NewApplicationError(mappings.QuizNotAvailableError, nil)
	case err != nil:
		return nil, apperrors.NewInternalError(err)
	}

	saved := make([]DraftAnswerData, 0)
	if result.Resumed {
		stored, err := app.Repositories.QuizAttempt.ListAnswers(ctx, result.Attempt.ID)
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		for _, a := range stored {
			saved = append(saved, DraftAnswerData{QuestionID: a.QuestionID, AnswerText: a.AnswerText})
		}
	}

	return &StartOutput{
		Attempt:      toAttemptData(result.Attempt),
		Questions:    toStudentQuestionData(questions),
		TimeLimit:    q.TimeLimitSecs,
		Resumed:      result.Resumed,
		SavedAnswers: saved,
		ServerNow:    time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	}, nil
}
