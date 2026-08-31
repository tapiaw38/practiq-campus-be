package course_group

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.CourseGroup) (string, error)
	Get(context.Context, string) (*domain.CourseGroup, error)
	ListByCourse(context.Context, string) ([]domain.CourseGroup, error)
	Delete(context.Context, string) error

	AddMember(ctx context.Context, groupID, userID string) error
	RemoveMember(ctx context.Context, groupID, userID string) error
	IsMember(ctx context.Context, groupID, userID string) (bool, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
