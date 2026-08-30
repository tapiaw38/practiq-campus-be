package course_group

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, g domain.CourseGroup) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `INSERT INTO course_groups (course_id, name) VALUES ($1, $2) RETURNING id`, g.CourseID, g.Name).Scan(&id)
	return id, err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM course_groups WHERE id=$1`, id)
	return err
}
