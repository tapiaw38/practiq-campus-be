package quiz_attempt

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.QuizAttempt) (string, error)
	Get(context.Context, string) (*domain.QuizAttempt, error)
	CountByUser(context.Context, string, string) (int, error)
	ListByQuiz(context.Context, string) ([]domain.QuizAttempt, error)
	ListMine(context.Context, string, string) ([]domain.QuizAttempt, error)
	Submit(context.Context, string, int, int) error

	SaveAnswers(context.Context, string, []domain.QuizAnswer) error
	ListAnswers(context.Context, string) ([]domain.QuizAnswer, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
