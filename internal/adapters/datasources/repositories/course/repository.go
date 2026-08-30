package course

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Course) (string, error)
	Get(context.Context, string) (*domain.Course, error)
	GetBySlug(context.Context, string) (*domain.Course, error)
	List(context.Context, ListFilter) ([]domain.Course, error)
	Update(context.Context, string, domain.Course) error
	Delete(context.Context, string) error
}

// ListFilter narrows a course listing. OwnerID lists a teacher's own
// courses; EnrolledUserID lists courses a given user is enrolled in
// (joins enrollments); neither set lists every published course.
type ListFilter struct {
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
