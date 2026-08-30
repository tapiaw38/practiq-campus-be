package domain

import "time"

const (
	SubmissionStatusSubmitted = "submitted"
	SubmissionStatusGraded    = "graded"
)

// Submission carries its own grade fields (Score/Feedback/GradedAt) rather
// than a separate grades table — in practice a submission has exactly one
// grade, so splitting them into two tables would only add a join with no
// real independent lifecycle to justify it.
type Submission struct {
	ID           string
	AssignmentID string
	UserID       string
	UserName     string
	Content      string
	Status       string
	Score        *int
	Feedback     string
	SubmittedAt  time.Time
	GradedAt     *time.Time
}
