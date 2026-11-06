package assignment

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

var ErrCourseNotInTenant = errors.New("course does not belong to this tenant")

func (r *repository) Create(ctx context.Context, a domain.Assignment) (string, error) {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"courses c", "c.id = $1", "c.tenant_id", 11)
	if err != nil {
		return "", err
	}
	query := `
		INSERT INTO assignments (course_id, section_id, title, description, due_at, max_score, weight, visible_group_id, unlock_after_type, unlock_after_id)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		WHERE ` + guard + `
		RETURNING id
	`
	var id string
	err = r.db.QueryRowContext(ctx, query,
		a.CourseID, a.SectionID, a.Title, a.Description, a.DueAt, a.MaxScore,
		a.Weight, a.VisibleGroupID, a.UnlockAfterType, a.UnlockAfterID, tenantID,
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrCourseNotInTenant
	}
	return id, err
}
