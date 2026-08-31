package assignment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, a domain.Assignment) (string, error) {
	query := `
		INSERT INTO assignments (course_id, section_id, title, description, due_at, max_score, weight, visible_group_id, unlock_after_type, unlock_after_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query,
		a.CourseID, a.SectionID, a.Title, a.Description, a.DueAt, a.MaxScore, a.Weight, a.VisibleGroupID, a.UnlockAfterType, a.UnlockAfterID,
	).Scan(&id)
	return id, err
}
