package course_group

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

const listQuery = `
	SELECT g.id, g.course_id, g.name, g.created_at,
		COALESCE(array_agg(m.user_id) FILTER (WHERE m.user_id IS NOT NULL), '{}')
	FROM course_groups g
	LEFT JOIN course_group_members m ON m.group_id = g.id
`

func (r *repository) Get(ctx context.Context, id string) (*domain.CourseGroup, error) {
	row := r.db.QueryRowContext(ctx, listQuery+` WHERE g.id = $1 GROUP BY g.id`, id)
	var g domain.CourseGroup
	err := row.Scan(&g.ID, &g.CourseID, &g.Name, &g.CreatedAt, pq.Array(&g.MemberIDs))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.CourseGroup, error) {
	rows, err := r.db.QueryContext(ctx, listQuery+` WHERE g.course_id = $1 GROUP BY g.id ORDER BY g.created_at ASC`, courseID)
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
