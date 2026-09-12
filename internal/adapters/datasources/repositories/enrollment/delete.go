package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) Delete(ctx context.Context, id string) error {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"courses c", "c.id = enrollments.course_id", "c.tenant_id", 2)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, "DELETE FROM enrollments WHERE id = $1 AND "+guard, id, tenantID)
	return err
}
