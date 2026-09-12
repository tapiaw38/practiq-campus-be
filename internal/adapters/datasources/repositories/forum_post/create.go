package forum_post

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// ErrThreadNotInTenant is returned when the thread belongs to another
// institution's course.
var ErrThreadNotInTenant = errors.New("thread does not belong to this tenant")

func (r *repository) Create(ctx context.Context, p domain.ForumPost) (string, error) {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"forum_threads t JOIN courses c ON c.id = t.course_id", "t.id = $1", "c.tenant_id", 5)
	if err != nil {
		return "", err
	}
	var id string
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO forum_posts (thread_id, parent_post_id, author_id, body)
		SELECT $1, $2, $3, $4 WHERE `+guard+`
		RETURNING id`, p.ThreadID, p.ParentID, p.AuthorID, p.Body, tenantID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrThreadNotInTenant
	}
	return id, err
}
