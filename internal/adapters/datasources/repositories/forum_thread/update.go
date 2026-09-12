package forum_thread

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) Update(ctx context.Context, id string, t domain.ForumThread) error {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"courses c", "c.id = forum_threads.course_id", "c.tenant_id", 4)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx,
		`UPDATE forum_threads SET title = $2, description = $3 WHERE id = $1 AND `+guard,
		id, t.Title, t.Description, tenantID)
	return err
}
