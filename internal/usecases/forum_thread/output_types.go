package forum_thread

import "github.com/tapiaw38/practiq-campus-be/internal/domain"

type ThreadData struct {
	ID          string `json:"id"`
	CourseID    string `json:"course_id"`
	AuthorID    string `json:"author_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func toThreadData(t domain.ForumThread) ThreadData {
	return ThreadData{
		ID:          t.ID,
		CourseID:    t.CourseID,
		AuthorID:    t.AuthorID,
		Title:       t.Title,
		Description: t.Description,
		CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toThreadDataList(threads []domain.ForumThread) []ThreadData {
	out := make([]ThreadData, 0, len(threads))
	for _, t := range threads {
		out = append(out, toThreadData(t))
	}
	return out
}
