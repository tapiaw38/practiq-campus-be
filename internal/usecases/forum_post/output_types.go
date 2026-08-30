package forum_post

import "github.com/tapiaw38/practiq-campus-be/internal/domain"

type PostData struct {
	ID         string  `json:"id"`
	ThreadID   string  `json:"thread_id"`
	ParentID   *string `json:"parent_post_id,omitempty"`
	AuthorID   string  `json:"author_id"`
	AuthorName string  `json:"author_name"`
	Body       string  `json:"body"`
	CreatedAt  string  `json:"created_at"`
}

func toPostData(p domain.ForumPost) PostData {
	return PostData{
		ID:         p.ID,
		ThreadID:   p.ThreadID,
		ParentID:   p.ParentID,
		AuthorID:   p.AuthorID,
		AuthorName: p.AuthorName,
		Body:       p.Body,
		CreatedAt:  p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toPostDataList(posts []domain.ForumPost) []PostData {
	out := make([]PostData, 0, len(posts))
	for _, p := range posts {
		out = append(out, toPostData(p))
	}
	return out
}
