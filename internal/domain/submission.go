package domain

import "time"

const (
	SubmissionStatusSubmitted = "submitted"
	SubmissionStatusGraded    = "graded"
)

type Submission struct {
	ID           string
	AssignmentID string
	UserID       string
	UserName     string
	Content      string
	Status       string

	Version     int
	Score       *int
	Feedback    string
	SubmittedAt time.Time
	GradedAt    *time.Time
}
