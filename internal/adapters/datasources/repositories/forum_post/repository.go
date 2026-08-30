package forum_post

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.ForumPost) (string, error)
	Get(context.Context, string) (*domain.ForumPost, error)
	ListByThread(context.Context, string, ListOptions) ([]domain.ForumPost, error)
	CountRootsByThread(context.Context, string) (int, error)
}

type ListOptions struct{ Limit, Offset int }

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
