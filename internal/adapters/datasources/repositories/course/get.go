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

// Get reads one course of the tenant in context. A course id from another
// institution simply is not found: it is not a permission the caller is denied
// afterwards, it is a row the query never returns.
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

// GetBySlug takes the tenant explicitly because slugs are only unique within
// one: the same slug in two institutions is two different courses.
func (r *repository) GetBySlug(ctx context.Context, tenantID, slug string) (*domain.Course, error) {
	row := r.db.QueryRowContext(ctx, "SELECT "+selectCourseColumns+" FROM courses WHERE tenant_id = $1 AND slug = $2", tenantID, slug)
	return scanCourse(row)
}
