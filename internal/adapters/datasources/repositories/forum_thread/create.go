package forum_thread

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, t domain.ForumThread) (string, error) {
	query := `
		INSERT INTO forum_threads (course_id, author_id, title, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, t.CourseID, t.AuthorID, t.Title, t.Description).Scan(&id)
	return id, err
}
