package course_section

import (
	"context"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, s domain.CourseSection) error {
	_, e := r.db.ExecContext(ctx, `UPDATE course_sections SET title=$1, position=$2, updated_at=now() WHERE id=$3`, s.Title, s.Position, id)
	return e
}
func (r *repository) Delete(ctx context.Context, id string) error {
	_, e := r.db.ExecContext(ctx, `DELETE FROM course_sections WHERE id=$1`, id)
	return e
}
