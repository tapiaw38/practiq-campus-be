package course

import (
	"context"
	"github.com/lib/pq"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, c domain.Course) error {
	query := `
		UPDATE courses SET
			title = $1, description = $2, status = $3, start_date = $4, end_date = $5,
			practiq_subject_id = $6, labels = $7, updated_at = NOW()
		WHERE id = $8
	`
	_, err := r.db.ExecContext(ctx, query, c.Title, c.Description, c.Status, c.StartDate, c.EndDate, c.PractiqSubjectID, pq.Array(c.Labels), id)
	return err
}
