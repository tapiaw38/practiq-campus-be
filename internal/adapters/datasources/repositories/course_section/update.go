package course_section

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) Update(ctx context.Context, id string, s domain.CourseSection) error {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"courses c", "c.id = course_sections.course_id", "c.tenant_id", 5)
	if err != nil {
		return err
	}
	_, e := r.db.ExecContext(ctx,
		`UPDATE course_sections SET title=$1, description=$2, position=$3, updated_at=now() WHERE id=$4 AND `+guard,
		s.Title, s.Description, s.Position, id, tenantID)
	return e
}

func (r *repository) Delete(ctx context.Context, id string) error {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"courses c", "c.id = course_sections.course_id", "c.tenant_id", 2)
	if err != nil {
		return err
	}
	_, e := r.db.ExecContext(ctx, `DELETE FROM course_sections WHERE id=$1 AND `+guard, id, tenantID)
	return e
}
