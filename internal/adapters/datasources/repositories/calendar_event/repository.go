package calendar_event

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.CalendarEvent) (string, error)
	Get(context.Context, string) (*domain.CalendarEvent, error)
	ListByOwner(context.Context, string) ([]domain.CalendarEvent, error)
	Update(context.Context, string, domain.CalendarEvent) error
	Delete(context.Context, string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
