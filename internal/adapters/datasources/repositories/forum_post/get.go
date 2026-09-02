package forum_post

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.ForumPost, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, thread_id, parent_post_id, author_id, body, created_at FROM forum_posts WHERE id = $1", id)

	var p domain.ForumPost
	err := row.Scan(&p.ID, &p.ThreadID, &p.ParentID, &p.AuthorID, &p.Body, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}
