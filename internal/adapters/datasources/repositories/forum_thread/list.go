package forum_thread

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.ForumThread, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("courses", "c", "t.course_id").
		Where("t.course_id = ?", courseID).
		SQL(selectQualifiedThreadColumns, "forum_threads t")
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query+" ORDER BY t.created_at DESC", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []domain.ForumThread
	for rows.Next() {
		t, err := scanThread(rows)
		if err != nil {
			return nil, err
		}
		threads = append(threads, *t)
	}
	return threads, rows.Err()
}
