package course_section

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.CourseSection) (string, error)
	Get(context.Context, string) (*domain.CourseSection, error)
	ListByCourse(context.Context, string) ([]domain.CourseSection, error)
	Update(context.Context, string, domain.CourseSection) error
	Delete(context.Context, string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
