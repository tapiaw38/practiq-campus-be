package forum_thread

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// Qualified for the tenant join, which brings courses into scope.
const selectQualifiedThreadColumns = `t.id, t.course_id, t.author_id, t.title, t.description, t.created_at`

func scanThread(row interface{ Scan(...any) error }) (*domain.ForumThread, error) {
	var t domain.ForumThread
	err := row.Scan(&t.ID, &t.CourseID, &t.AuthorID, &t.Title, &t.Description, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.ForumThread, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "t.course_id").
		Where("t.id = ?", id).
		SQL(selectQualifiedThreadColumns, "forum_threads t")
	if err != nil {
		return nil, err
	}
	return scanThread(r.db.QueryRowContext(ctx, query, args...))
}
