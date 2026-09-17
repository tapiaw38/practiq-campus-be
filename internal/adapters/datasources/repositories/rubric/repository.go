package rubric

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	List(context.Context, string) ([]domain.RubricCriterion, error)
	Replace(context.Context, string, []domain.RubricCriterion) error
	Grade(ctx context.Context, submissionID string, total int, feedback string, scores []domain.RubricScore, expectedVersion int) error
	Scores(context.Context, string) ([]domain.RubricScore, error)
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db} }
