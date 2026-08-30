package forum_post

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
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
	query := `WITH RECURSIVE root_page AS (
		SELECT id, created_at AS root_created_at
		FROM forum_posts
		WHERE thread_id = $1 AND parent_post_id IS NULL
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
	SELECT p.id, p.thread_id, p.parent_post_id, p.author_id, COALESCE(cp.full_name, p.author_id), p.body, p.created_at
	FROM post_tree p
	LEFT JOIN campus_profiles cp ON cp.id = p.author_id
	ORDER BY p.root_created_at ASC, p.created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, threadID, limit, options.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.ForumPost
	for rows.Next() {
		var p domain.ForumPost
		if err := rows.Scan(&p.ID, &p.ThreadID, &p.ParentID, &p.AuthorID, &p.AuthorName, &p.Body, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *repository) CountRootsByThread(ctx context.Context, threadID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM forum_posts WHERE thread_id = $1 AND parent_post_id IS NULL`, threadID).Scan(&count)
	return count, err
}
