package forum_post

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) ListByThread(ctx context.Context, threadID string, options ListOptions) ([]domain.ForumPost, error) {
	limit := options.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if options.Offset < 0 {
		options.Offset = 0
	}
	// Pagination deliberately applies only to root posts. Every reply beneath
	// those roots travels with its conversation, so a thread is never split
	// across pages or counted as a general forum message.
	tenantID := tenantcontext.ID(ctx)
	if tenantID == "" {
		return nil, tenantcontext.ErrNoTenant
	}
	// Only the anchor needs the tenant: every other row in the tree descends
	// from a root of this thread by parent_post_id, so a thread in another
	// institution yields no roots and therefore no replies.
	query := `WITH RECURSIVE root_page AS (
		SELECT p.id, p.created_at AS root_created_at
		FROM forum_posts p
		JOIN forum_threads t ON t.id = p.thread_id
		JOIN courses c ON c.id = t.course_id AND c.tenant_id = $4
		WHERE p.thread_id = $1 AND p.parent_post_id IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	), post_tree AS (
		SELECT p.id, p.thread_id, p.parent_post_id, p.author_id, p.body, p.created_at, r.root_created_at
		FROM forum_posts p
		JOIN root_page r ON r.id = p.id
		UNION ALL
		SELECT child.id, child.thread_id, child.parent_post_id, child.author_id, child.body, child.created_at, tree.root_created_at
		FROM forum_posts child
		JOIN post_tree tree ON child.parent_post_id = tree.id
	)
	SELECT id, thread_id, parent_post_id, author_id, body, created_at
	FROM post_tree
	ORDER BY root_created_at ASC, created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, threadID, limit, options.Offset, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.ForumPost
	for rows.Next() {
		var p domain.ForumPost
		if err := rows.Scan(&p.ID, &p.ThreadID, &p.ParentID, &p.AuthorID, &p.Body, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *repository) CountRootsByThread(ctx context.Context, threadID string) (int, error) {
	tenantID := tenantcontext.ID(ctx)
	if tenantID == "" {
		return 0, tenantcontext.ErrNoTenant
	}
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM forum_posts p
		JOIN forum_threads t ON t.id = p.thread_id
		JOIN courses c ON c.id = t.course_id AND c.tenant_id = $2
		WHERE p.thread_id = $1 AND p.parent_post_id IS NULL`, threadID, tenantID).Scan(&count)
	return count, err
}
