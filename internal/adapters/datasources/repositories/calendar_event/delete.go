package calendar_event

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) Delete(ctx context.Context, id string) error {
	tenantID := tenantcontext.ID(ctx)
	if tenantID == "" {
		return tenantcontext.ErrNoTenant
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM calendar_events WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	return err
}
