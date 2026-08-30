package assignment

import (
	"context"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, a domain.Assignment) error {
	_, e := r.db.ExecContext(ctx, `UPDATE assignments SET section_id=$1,title=$2,description=$3,due_at=$4,max_score=$5,updated_at=now() WHERE id=$6`, a.SectionID, a.Title, a.Description, a.DueAt, a.MaxScore, id)
	return e
}
func (r *repository) Delete(ctx context.Context, id string) error {
	_, e := r.db.ExecContext(ctx, `DELETE FROM assignments WHERE id=$1`, id)
	return e
}
