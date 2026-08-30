package domain

import "time"

type Assignment struct {
	ID          string
	CourseID    string
	SectionID   *string
	Title       string
	Description string
	DueAt       *time.Time
	MaxScore    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
