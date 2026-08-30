package assignment

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Assignment) (string, error)
	Get(context.Context, string) (*domain.Assignment, error)
	ListByCourse(context.Context, string) ([]domain.Assignment, error)
	Update(context.Context, string, domain.Assignment) error
	Delete(context.Context, string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
