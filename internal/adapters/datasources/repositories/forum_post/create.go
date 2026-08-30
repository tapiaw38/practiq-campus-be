package forum_post

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, p domain.ForumPost) (string, error) {
	query := `
		INSERT INTO forum_posts (thread_id, parent_post_id, author_id, body)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, p.ThreadID, p.ParentID, p.AuthorID, p.Body).Scan(&id)
	return id, err
}
