package course

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func scanCourse(row *sql.Row) (*domain.Course, error) {
	var c domain.Course
	err := row.Scan(
		&c.ID, &c.TenantID, &c.OwnerID, &c.Title, &c.Slug, &c.Description, &c.Status,
		&c.StartDate, &c.EndDate, &c.CreatedAt, &c.UpdatedAt, &c.PractiqSubjectID, pq.Array(&c.Labels),
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

const selectCourseColumns = `
	id, tenant_id, owner_id, title, slug, description, status, start_date, end_date, created_at, updated_at, practiq_subject_id, labels
`

func (r *repository) Get(ctx context.Context, id string) (*domain.Course, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Where("c.id = ?", id).
		Own("c.tenant_id").
		SQL(selectQualifiedCourseColumns, "courses c")
	if err != nil {
		return nil, err
	}
	return scanCourse(r.db.QueryRowContext(ctx, query, args...))
}

func (r *repository) GetBySlug(ctx context.Context, tenantID, slug string) (*domain.Course, error) {

	if tenantID == "" {
		tenantID = tenantcontext.ID(ctx)
	}
	if tenantID == "" {
		return nil, tenantcontext.ErrNoTenant
	}
	row := r.db.QueryRowContext(ctx, "SELECT "+selectCourseColumns+" FROM courses WHERE tenant_id = $1 AND slug = $2", tenantID, slug)
	return scanCourse(row)
}
