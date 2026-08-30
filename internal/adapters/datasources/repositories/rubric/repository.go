package rubric

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	List(context.Context, string) ([]domain.RubricCriterion, error)
	Replace(context.Context, string, []domain.RubricCriterion) error
	Grade(context.Context, string, int, string, []domain.RubricScore) error
	Scores(context.Context, string) ([]domain.RubricScore, error)
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db} }
