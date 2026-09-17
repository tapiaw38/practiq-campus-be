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
	Content string
	Status  string
	// Version increments on every resubmission. A grade carries the version
	// it was written against, so a teacher who reviewed one version cannot
	// have their grade applied to a newer one the student sent meanwhile.
	Version     int
	Score       *int
	Feedback    string
	SubmittedAt time.Time
	GradedAt    *time.Time
}
