package forum_thread

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.ForumThread, error) {
	query := `SELECT ` + selectThreadColumns + ` FROM forum_threads WHERE course_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []domain.ForumThread
	for rows.Next() {
		var t domain.ForumThread
		if err := rows.Scan(&t.ID, &t.CourseID, &t.AuthorID, &t.Title, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		threads = append(threads, t)
	}
	return threads, rows.Err()
}
