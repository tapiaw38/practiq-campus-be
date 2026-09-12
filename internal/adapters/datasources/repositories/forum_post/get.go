package forum_post

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// A post reaches its tenant through two hops: thread, then course. Neither
// forum table carries a tenant of its own.
func (r *repository) Get(ctx context.Context, id string) (*domain.ForumPost, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Join("forum_threads", "t", "p.thread_id").
		Through("courses", "c", "t.course_id").
		Where("p.id = ?", id).
		SQL("p.id, p.thread_id, p.parent_post_id, p.author_id, p.body, p.created_at", "forum_posts p")
	if err != nil {
		return nil, err
	}
	var p domain.ForumPost
	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&p.ID, &p.ThreadID, &p.ParentID, &p.AuthorID, &p.Body, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}
