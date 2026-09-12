package assignment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// Update and Delete carry the tenant guard in the same WHERE as the id. A
// statement naming a row in another institution matches nothing: there is no
// window in which it is changed and then corrected.
func (r *repository) Update(ctx context.Context, id string, a domain.Assignment) error {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"courses c", "c.id = assignments.course_id", "c.tenant_id", 11)
	if err != nil {
		return err
	}
	_, e := r.db.ExecContext(ctx, `
		UPDATE assignments SET section_id=$1,title=$2,description=$3,due_at=$4,max_score=$5,
		weight=$6,visible_group_id=$7,unlock_after_type=$8,unlock_after_id=$9,updated_at=now()
		WHERE id=$10 AND `+guard,
		a.SectionID, a.Title, a.Description, a.DueAt, a.MaxScore, a.Weight,
		a.VisibleGroupID, a.UnlockAfterType, a.UnlockAfterID, id, tenantID)
	return e
}

func (r *repository) Delete(ctx context.Context, id string) error {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"courses c", "c.id = assignments.course_id", "c.tenant_id", 2)
	if err != nil {
		return err
	}
	_, e := r.db.ExecContext(ctx, `DELETE FROM assignments WHERE id=$1 AND `+guard, id, tenantID)
	return e
}
