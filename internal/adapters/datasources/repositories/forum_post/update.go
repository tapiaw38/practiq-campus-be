package forum_post

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) Update(ctx context.Context, id, body string) error {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"forum_threads t JOIN courses c ON c.id = t.course_id", "t.id = forum_posts.thread_id", "c.tenant_id", 3)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `UPDATE forum_posts SET body = $2 WHERE id = $1 AND `+guard, id, body, tenantID)
	return err
}
