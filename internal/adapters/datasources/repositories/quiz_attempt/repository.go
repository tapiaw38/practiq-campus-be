package quiz_attempt

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Get(context.Context, string) (*domain.QuizAttempt, error)
	CountByUser(context.Context, string, string) (int, error)
	ListByQuiz(context.Context, string) ([]domain.QuizAttempt, error)
	ListMine(context.Context, string, string) ([]domain.QuizAttempt, error)

	// StartOrResume and SubmitGraded are the only ways in and out of an
	// attempt. The pieces they are built from — create, save answers, close —
	// are deliberately not exposed: run separately they raced each other,
	// which is what made reloading a page burn an attempt and a failed submit
	// strand answers against an attempt left open.
	StartOrResume(context.Context, StartParams) (*StartResult, error)
	SaveDraft(ctx context.Context, attemptID, userID string, answers []domain.QuizAnswer) error
	// FinalizeExpired settles an attempt whose deadline passed unsubmitted,
	// grading whatever was saved. Idempotent and a no-op on a live attempt,
	// so read paths can call it freely.
	FinalizeExpired(ctx context.Context, attemptID string, answers []domain.QuizAnswer, score, maxScore int) (*domain.QuizAttempt, error)
	SubmitGraded(ctx context.Context, attemptID, userID string, answers []domain.QuizAnswer, score, maxScore int) (*domain.QuizAttempt, error)

	ListAnswers(context.Context, string) ([]domain.QuizAnswer, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
