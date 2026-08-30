package quiz

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Quiz) (string, error)
	Get(context.Context, string) (*domain.Quiz, error)
	ListByCourse(context.Context, string) ([]domain.Quiz, error)
	Update(context.Context, string, domain.Quiz) error
	Delete(context.Context, string) error

	ListQuestions(context.Context, string) ([]domain.QuizQuestion, error)
	ReplaceQuestions(context.Context, string, []domain.QuizQuestion) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
