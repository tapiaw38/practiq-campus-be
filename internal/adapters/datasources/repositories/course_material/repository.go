package course_material

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.CourseMaterial) (string, error)
	Get(context.Context, string) (*domain.CourseMaterial, error)
	ListByCourse(context.Context, string) ([]domain.CourseMaterial, error)
	Delete(context.Context, string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
