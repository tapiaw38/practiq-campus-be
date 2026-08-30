package submission

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Submission) (string, error)
	Resubmit(context.Context, string, string) error
	Get(context.Context, string) (*domain.Submission, error)
	GetByAssignmentAndUser(context.Context, string, string) (*domain.Submission, error)
	ListByAssignment(context.Context, string) ([]domain.Submission, error)
	Grade(ctx context.Context, id string, score int, feedback string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
