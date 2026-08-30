package forum_post

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.ForumPost, error) {
	row := r.db.QueryRowContext(ctx, "SELECT p.id, p.thread_id, p.parent_post_id, p.author_id, COALESCE(cp.full_name, p.author_id), p.body, p.created_at FROM forum_posts p LEFT JOIN campus_profiles cp ON cp.id = p.author_id WHERE p.id = $1", id)

	var p domain.ForumPost
	err := row.Scan(&p.ID, &p.ThreadID, &p.ParentID, &p.AuthorID, &p.AuthorName, &p.Body, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}
