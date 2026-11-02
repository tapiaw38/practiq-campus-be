package domain

import "time"

const (
	UnlockAfterAssignment = "assignment"
	UnlockAfterQuiz       = "quiz"
)

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

	Weight int

	VisibleGroupID *string

	UnlockAfterType *string
	UnlockAfterID   *string
}
