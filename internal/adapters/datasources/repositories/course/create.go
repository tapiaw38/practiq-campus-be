package course

import (
	"context"
	"github.com/lib/pq"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) Create(ctx context.Context, c domain.Course) (string, error) {
	if c.TenantID == "" {
		c.TenantID = tenantcontext.ID(ctx)
	}
	// tenant_id is still nullable in the schema, so an insert with no tenant
	// would succeed and leave a course belonging to no institution — visible
	// to nobody, counted by nothing, and impossible to attribute later.
	if c.TenantID == "" {
		return "", tenantcontext.ErrNoTenant
	}
	query := `
		INSERT INTO courses (tenant_id, owner_id, title, slug, description, status, start_date, end_date, practiq_subject_id, labels)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query,
		c.TenantID, c.OwnerID, c.Title, c.Slug, c.Description, c.Status, c.StartDate, c.EndDate, c.PractiqSubjectID, pq.Array(c.Labels),
	).Scan(&id)
	return id, err
}
