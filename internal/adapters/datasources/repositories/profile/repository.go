package profile

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Upsert(context.Context, domain.Profile) error
	Get(context.Context, string) (*domain.Profile, error)
	GetByEmail(context.Context, string) (*domain.Profile, error)
	ListByType(context.Context, string) ([]domain.Profile, error)
	ListAll(context.Context, string, int, int) ([]domain.Profile, int, error)
	SetBlocked(context.Context, string, bool) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
