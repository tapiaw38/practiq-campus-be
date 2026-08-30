package course

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, c domain.Course) error {
	query := `
		UPDATE courses SET
			title = $1, description = $2, status = $3, start_date = $4, end_date = $5, updated_at = NOW()
		WHERE id = $6
	`
	_, err := r.db.ExecContext(ctx, query, c.Title, c.Description, c.Status, c.StartDate, c.EndDate, id)
	return err
}
