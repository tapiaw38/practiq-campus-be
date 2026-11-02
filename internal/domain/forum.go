package domain

import "time"

type ForumThread struct {
	ID          string
	CourseID    string
	AuthorID    string
	AuthorName  string
	Title       string
	Description string
	CreatedAt   time.Time
}

type ForumPost struct {
	ID         string
	ThreadID   string
	ParentID   *string
	AuthorID   string
	AuthorName string
	Body       string
	CreatedAt  time.Time
	DeletedAt  *time.Time
}
