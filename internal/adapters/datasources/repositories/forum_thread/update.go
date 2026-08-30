package forum_thread

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, t domain.ForumThread) error {
	_, err := r.db.ExecContext(ctx, `UPDATE forum_threads SET title = $2, description = $3 WHERE id = $1`, id, t.Title, t.Description)
	return err
}
