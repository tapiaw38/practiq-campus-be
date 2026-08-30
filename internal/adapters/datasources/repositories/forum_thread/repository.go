package forum_thread

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.ForumThread) (string, error)
	Get(context.Context, string) (*domain.ForumThread, error)
	ListByCourse(context.Context, string) ([]domain.ForumThread, error)
	Update(context.Context, string, domain.ForumThread) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
