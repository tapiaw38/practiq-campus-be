package course

import (
	"context"

	"github.com/lib/pq"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// Qualified with c.: the JOIN below brings in enrollments, which also has
// id/status columns, so the bare names in selectCourseColumns (fine for the
// single-table Get/GetBySlug queries) are ambiguous here.
const selectQualifiedCourseColumns = `
	c.id, c.tenant_id, c.owner_id, c.title, c.slug, c.description, c.status, c.start_date, c.end_date, c.created_at, c.updated_at, c.practiq_subject_id, c.labels
`

func (r *repository) List(ctx context.Context, filter ListFilter) ([]domain.Course, error) {
	// No tenant is refused rather than left unfiltered. This listing used to
	// skip the condition when none was set, which answered a request that
	// never passed RequireTenant with every institution's courses.
	q := tenantcontext.NewQuery(ctx).Own("c.tenant_id")

	from := "courses c"
	if filter.EnrolledUserID != "" {
		from += " JOIN enrollments e ON e.course_id = c.id"
		q.Where("e.user_id = ? AND e.status = 'active'", filter.EnrolledUserID)
	}
	if filter.OwnerID != "" {
		q.Where("c.owner_id = ?", filter.OwnerID)
	}
	if filter.PublishedOnly {
		q.Where("c.status = 'published'")
	}

	query, args, err := q.SQL(selectQualifiedCourseColumns, from)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query+" ORDER BY c.created_at DESC", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []domain.Course
	for rows.Next() {
		var c domain.Course
		if err := rows.Scan(
			&c.ID, &c.TenantID, &c.OwnerID, &c.Title, &c.Slug, &c.Description, &c.Status,
			&c.StartDate, &c.EndDate, &c.CreatedAt, &c.UpdatedAt, &c.PractiqSubjectID, pq.Array(&c.Labels),
		); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, rows.Err()
}
