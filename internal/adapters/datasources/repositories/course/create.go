package course

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, c domain.Course) (string, error) {
	query := `
		INSERT INTO courses (owner_id, title, slug, description, status, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query,
		c.OwnerID, c.Title, c.Slug, c.Description, c.Status, c.StartDate, c.EndDate,
	).Scan(&id)
	return id, err
}
