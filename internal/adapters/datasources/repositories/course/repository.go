package course

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Course) (string, error)
	Get(context.Context, string) (*domain.Course, error)
	GetBySlug(context.Context, string, string) (*domain.Course, error)
	List(context.Context, ListFilter) ([]domain.Course, error)
	Update(context.Context, string, domain.Course) error
	Delete(context.Context, string) error
}

type ListFilter struct {
	TenantID       string
	OwnerID        string
	EnrolledUserID string
	PublishedOnly  bool
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
