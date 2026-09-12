package course_group

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// The group listing aggregates its members, and joins courses so the tenant
// constrains which groups exist at all.
const groupColumns = `g.id, g.course_id, g.name, g.created_at,
	COALESCE(array_agg(m.user_id) FILTER (WHERE m.user_id IS NOT NULL), '{}')`

const groupFrom = `course_groups g
	LEFT JOIN course_group_members m ON m.group_id = g.id`

func (r *repository) Get(ctx context.Context, id string) (*domain.CourseGroup, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "g.course_id").
		Where("g.id = ?", id).
		SQL(groupColumns, groupFrom)
	if err != nil {
		return nil, err
	}
	var g domain.CourseGroup
	err = r.db.QueryRowContext(ctx, query+" GROUP BY g.id", args...).
		Scan(&g.ID, &g.CourseID, &g.Name, &g.CreatedAt, pq.Array(&g.MemberIDs))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.CourseGroup, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "g.course_id").
		Where("g.course_id = ?", courseID).
		SQL(groupColumns, groupFrom)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query+" GROUP BY g.id ORDER BY g.created_at ASC", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]domain.CourseGroup, 0)
	for rows.Next() {
		var g domain.CourseGroup
		if err := rows.Scan(&g.ID, &g.CourseID, &g.Name, &g.CreatedAt, pq.Array(&g.MemberIDs)); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}
