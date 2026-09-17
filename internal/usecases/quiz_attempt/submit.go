package quiz_attempt

import (
	"context"
	"errors"
	"time"

	attemptrepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/quiz_attempt"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	SubmitUsecase interface {
		Execute(ctx context.Context, requesterID, attemptID string, answers []AnswerInput) (*SubmitOutput, apperrors.ApplicationError)
	}

	submitUsecase struct {
		contextFactory appcontext.Factory
	}

	AnswerInput struct {
		QuestionID string
		AnswerText string
	}

	SubmitOutput struct {
		Attempt AttemptData        `json:"attempt"`
		Results []AnswerResultData `json:"results"`
	}
)

func NewSubmitUsecase(contextFactory appcontext.Factory) SubmitUsecase {
	return &submitUsecase{contextFactory: contextFactory}
}

// Execute grades the attempt itself: every question type supported here is
// decided by exact comparison (domain.GradeQuizAnswer), never by an external
// call, so a result is available the instant the student submits.
func (u *submitUsecase) Execute(ctx context.Context, requesterID, attemptID string, answers []AnswerInput) (*SubmitOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	attempt, err := app.Repositories.QuizAttempt.Get(ctx, attemptID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	if attempt == nil {
		return nil, apperrors.NewApplicationError(mappings.QuizAttemptNotFoundError, nil)
	}
	if attempt.UserID != requesterID {
		return nil, apperrors.NewForbiddenError()
	}
	if attempt.SubmittedAt != nil {
		return nil, apperrors.NewApplicationError(mappings.QuizAttemptAlreadySubmittedError, nil)
	}
	// Cheap early exit; the deadline is enforced for real under the row lock
	// in SubmitGraded, since this read and that write are not the same moment.
	if attempt.Expired(time.Now()) {
		return nil, apperrors.NewApplicationError(mappings.QuizAttemptExpiredError, nil)
	}

	q, err := app.Repositories.Quiz.Get(ctx, attempt.QuizID)
	if err != nil || q == nil {
		return nil, apperrors.NewInternalError(err)
	}

	// Handing work in is an act of participation, so it asks the same
	// question starting does: are you still in this course, and still in the
	// group this quiz is for? Only ownership of the attempt was checked here,
	// so a student removed from the course mid-exam could still submit and
	// have a grade recorded against a course they no longer belong to — while
	// being refused if they tried to start the very same quiz.
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

	questions, err := app.Repositories.Quiz.ListQuestions(ctx, attempt.QuizID)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}

	byQuestion := make(map[string]string, len(answers))
	for _, a := range answers {
		byQuestion[a.QuestionID] = a.AnswerText
	}

	score, maxScore := 0, 0
	quizAnswers := make([]domain.QuizAnswer, 0, len(questions))
	results := make([]AnswerResultData, 0, len(questions))
	for _, question := range questions {
		maxScore += question.Points
		answerText := byQuestion[question.ID]
		isCorrect := domain.GradeQuizAnswer(question, answerText)
		if isCorrect {
			score += question.Points
		}
		quizAnswers = append(quizAnswers, domain.QuizAnswer{AttemptID: attemptID, QuestionID: question.ID, AnswerText: answerText, IsCorrect: isCorrect})
		results = append(results, AnswerResultData{
			QuestionID:    question.ID,
			Statement:     question.Statement,
			AnswerText:    answerText,
			CorrectAnswer: question.CorrectAnswer,
			IsCorrect:     isCorrect,
			Points:        question.Points,
		})
	}

	// Answers and the closing of the attempt go in together or not at all:
	// separately, a failure between them left answers saved against an
	// attempt still marked open, and two submissions could interleave into a
	// grade that did not match the answers stored.
	updated, err := app.Repositories.QuizAttempt.SubmitGraded(ctx, attemptID, requesterID, quizAnswers, score, maxScore)
	switch {
	case errors.Is(err, attemptrepo.ErrAttemptAlreadySubmitted):
		return nil, apperrors.NewApplicationError(mappings.QuizAttemptAlreadySubmittedError, nil)
	case errors.Is(err, attemptrepo.ErrAttemptExpired):
		// Ran out of time while waiting for the lock or being graded. Nothing
		// was stored, and saying "already submitted" would have the student
		// waiting for a grade that is never coming.
		return nil, apperrors.NewApplicationError(mappings.QuizAttemptExpiredError, nil)
	case errors.Is(err, attemptrepo.ErrNotYourAttempt):
		return nil, apperrors.NewForbiddenError()
	case err != nil:
		return nil, apperrors.NewInternalError(err)
	}

	// Only a submission that actually closed the attempt gets here, so the
	// teacher is notified once per attempt rather than once per request.
	if course, err := app.Repositories.Course.Get(ctx, q.CourseID); err == nil && course != nil {
		_ = app.Repositories.Notification.Create(ctx, domain.Notification{
			UserID: course.OwnerID,
			Type:   "quiz_submitted",
			Title:  "Entrega de evaluación: " + q.Title,
			Body:   "Un alumno completó la evaluación",
			Data:   `{"quiz_id":"` + q.ID + `","attempt_id":"` + attemptID + `","course_id":"` + q.CourseID + `"}`,
		})
	}

	return &SubmitOutput{Attempt: toAttemptData(*updated), Results: results}, nil
}
