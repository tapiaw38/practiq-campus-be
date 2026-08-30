package forum_thread

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

const selectThreadColumns = `id, course_id, author_id, title, description, created_at`

func (r *repository) Get(ctx context.Context, id string) (*domain.ForumThread, error) {
	row := r.db.QueryRowContext(ctx, "SELECT "+selectThreadColumns+" FROM forum_threads WHERE id = $1", id)

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
