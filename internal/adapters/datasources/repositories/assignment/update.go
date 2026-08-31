package assignment

import (
	"context"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, a domain.Assignment) error {
	_, e := r.db.ExecContext(ctx, `UPDATE assignments SET section_id=$1,title=$2,description=$3,due_at=$4,max_score=$5,weight=$6,visible_group_id=$7,unlock_after_type=$8,unlock_after_id=$9,updated_at=now() WHERE id=$10`, a.SectionID, a.Title, a.Description, a.DueAt, a.MaxScore, a.Weight, a.VisibleGroupID, a.UnlockAfterType, a.UnlockAfterID, id)
	return e
}
func (r *repository) Delete(ctx context.Context, id string) error {
	_, e := r.db.ExecContext(ctx, `DELETE FROM assignments WHERE id=$1`, id)
	return e
}
