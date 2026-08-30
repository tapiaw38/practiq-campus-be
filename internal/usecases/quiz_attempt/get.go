package quiz_attempt

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	GetUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, attemptID string) (*GetOutput, apperrors.ApplicationError)
	}

	getUsecase struct {
		contextFactory appcontext.Factory
	}

	GetOutput struct {
		Attempt AttemptData        `json:"attempt"`
		Results []AnswerResultData `json:"results"`
	}
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{contextFactory: contextFactory}
}

// Execute lets an attempt's own student, the course owner or a superadmin
// review it. Correct answers are only ever included once the attempt was
// submitted — reading an in-progress attempt would hand the answer key to
// whoever is still taking it.
func (u *getUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, attemptID string) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	attempt, err := app.Repositories.QuizAttempt.Get(ctx, attemptID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if attempt == nil {
		return nil, apperrors.NewApplicationError(mappings.QuizAttemptNotFoundError, nil)
	}

	q, err := app.Repositories.Quiz.Get(ctx, attempt.QuizID)
	if err != nil || q == nil {
		return nil, apperrors.NewInternalError(err)
	}

	if attempt.UserID != requesterID && !isSuperAdmin {
		course, err := app.Repositories.Course.Get(ctx, q.CourseID)
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		if course == nil || course.OwnerID != requesterID {
			return nil, apperrors.NewForbiddenError()
		}
	}

	results := make([]AnswerResultData, 0)
	if attempt.SubmittedAt != nil {
		questions, err := app.Repositories.Quiz.ListQuestions(ctx, attempt.QuizID)
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		byID := make(map[string]string, len(questions))
		points := make(map[string]int, len(questions))
		statements := make(map[string]string, len(questions))
		for _, question := range questions {
			byID[question.ID] = question.CorrectAnswer
			points[question.ID] = question.Points
			statements[question.ID] = question.Statement
		}
		answers, err := app.Repositories.QuizAttempt.ListAnswers(ctx, attemptID)
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		for _, a := range answers {
			results = append(results, AnswerResultData{
				QuestionID:    a.QuestionID,
				Statement:     statements[a.QuestionID],
				AnswerText:    a.AnswerText,
				CorrectAnswer: byID[a.QuestionID],
				IsCorrect:     a.IsCorrect,
				Points:        points[a.QuestionID],
			})
		}
	}

	return &GetOutput{Attempt: toAttemptData(*attempt), Results: results}, nil
}
