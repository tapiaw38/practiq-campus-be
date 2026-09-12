package tenant

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Tenant) (*domain.Tenant, error)
	Get(context.Context, string) (*domain.Tenant, error)
	GetBySchoolID(context.Context, string) (*domain.Tenant, error)
	List(context.Context) ([]domain.Tenant, error)
	SetStatus(context.Context, string, string) (*domain.Tenant, error)
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }
