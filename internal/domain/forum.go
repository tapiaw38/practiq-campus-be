package domain

import "time"

// A course has one flat list of threads directly (no separate "forum"
// entity per course) — nobody asked for multiple named forums per course,
// and Moodle-style single-board-per-course covers the actual use case.
type ForumThread struct {
	ID          string
	CourseID    string
	AuthorID    string
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
}
