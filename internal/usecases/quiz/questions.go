package quiz

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

// QuestionInput is what the teacher sends when replacing a quiz's questions.
type QuestionInput struct {
	Type          string
	Statement     string
	Options       []string
	CorrectAnswer string
	Points        int
}

type QuestionsUsecase interface {
	// List is the teacher-facing view, correct answers included.
	List(ctx context.Context, quizID, requesterID string, isSuperAdmin bool) ([]QuestionData, apperrors.ApplicationError)
	Replace(ctx context.Context, quizID, requesterID string, isSuperAdmin bool, in []QuestionInput) apperrors.ApplicationError
}

type questionsUsecase struct{ contextFactory appcontext.Factory }

func NewQuestionsUsecase(contextFactory appcontext.Factory) QuestionsUsecase {
	return &questionsUsecase{contextFactory: contextFactory}
}

func (u *questionsUsecase) quizOwnedByRequester(ctx context.Context, quizID, requesterID string, isSuperAdmin bool) (*domain.Quiz, apperrors.ApplicationError) {
	app := u.contextFactory()
	q, err := app.Repositories.Quiz.Get(ctx, quizID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.QuizGetError, err)
	}
	if q == nil {
		return nil, apperrors.NewApplicationError(mappings.QuizNotFoundError, nil)
	}
	if appErr := requesterOwnsCourse(ctx, app, requesterID, isSuperAdmin, q.CourseID); appErr != nil {
		return nil, appErr
	}
	return q, nil
}

func (u *questionsUsecase) List(ctx context.Context, quizID, requesterID string, isSuperAdmin bool) ([]QuestionData, apperrors.ApplicationError) {
	if _, appErr := u.quizOwnedByRequester(ctx, quizID, requesterID, isSuperAdmin); appErr != nil {
		return nil, appErr
	}
	questions, err := u.contextFactory().Repositories.Quiz.ListQuestions(ctx, quizID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return toQuestionDataList(questions), nil
}

func (u *questionsUsecase) Replace(ctx context.Context, quizID, requesterID string, isSuperAdmin bool, in []QuestionInput) apperrors.ApplicationError {
	if _, appErr := u.quizOwnedByRequester(ctx, quizID, requesterID, isSuperAdmin); appErr != nil {
		return appErr
	}

	questions := make([]domain.QuizQuestion, 0, len(in))
	for _, q := range in {
		points := q.Points
		if points <= 0 {
			points = 1
		}
		question := domain.QuizQuestion{
			QuizID:        quizID,
			Type:          q.Type,
			Statement:     q.Statement,
			Options:       q.Options,
			CorrectAnswer: q.CorrectAnswer,
			Points:        points,
		}
		if err := domain.ValidateQuizQuestion(question); err != nil {
			return apperrors.NewBadRequestError(err.Error())
		}
		questions = append(questions, question)
	}

	app := u.contextFactory()

	// Once a student has attempted the quiz, its questions are locked: editing
	// them would retroactively change what an already-graded attempt was
	// judged against.
	attempts, err := app.Repositories.QuizAttempt.ListByQuiz(ctx, quizID)
	if err != nil {
		return apperrors.NewInternalError(err)
	}
	if len(attempts) > 0 {
		return apperrors.NewBadRequestError("questions cannot be changed after a student has attempted the quiz")
	}

	if err = app.Repositories.Quiz.ReplaceQuestions(ctx, quizID, questions); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}
