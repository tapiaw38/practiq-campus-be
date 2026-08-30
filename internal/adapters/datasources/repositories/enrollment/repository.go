package enrollment

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Enrollment) (string, error)
	Get(context.Context, string) (*domain.Enrollment, error)
	GetByCourseAndUser(context.Context, string, string) (*domain.Enrollment, error)
	ListByCourse(context.Context, string) ([]domain.Enrollment, error)
	ListByUser(context.Context, string) ([]domain.Enrollment, error)
	Delete(context.Context, string) error
	// SharesCourseWith is true when userA and userB are on the same side of
	// at least one course together — one owns it and the other is enrolled,
	// or both are enrolled — regardless of which is which.
	SharesCourseWith(ctx context.Context, userA, userB string) (bool, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
