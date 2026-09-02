package profile

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Upsert(context.Context, domain.Profile) error
	Get(context.Context, string) (*domain.Profile, error)
	ListByType(context.Context, string) ([]domain.Profile, error)
	// ListAll has no text search — name/email no longer live locally, so
	// search-by-text now happens app-side against auth-api-be identity.
	// See usecases/profile.ListStudentsUsecase.
	ListAll(context.Context, int, int) ([]domain.Profile, int, error)
	SetBlocked(context.Context, string, bool) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
