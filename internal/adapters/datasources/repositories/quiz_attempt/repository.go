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

	StartOrResume(context.Context, StartParams) (*StartResult, error)
	SaveDraft(ctx context.Context, attemptID, userID string, answers []domain.QuizAnswer) error

	FinalizeExpired(ctx context.Context, attemptID string) (*domain.QuizAttempt, error)
	SubmitGraded(ctx context.Context, attemptID, userID string, answers []domain.QuizAnswer, score, maxScore int) (*domain.QuizAttempt, error)

	ListAnswers(context.Context, string) ([]domain.QuizAnswer, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
