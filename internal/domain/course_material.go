package domain

import "time"

const (
	MaterialKindFile = "file"
	MaterialKindLink = "link"
)

// CourseMaterial is a single piece of course content: either an uploaded file
// (URL points at the private bucket, presigned at read time) or an external
// link the teacher pasted. One field holds both because from the reader's side
// they are the same thing — something to open.
type CourseMaterial struct {
	ID           string
	CourseID     string
	AssignmentID *string
	SectionID    *string
	UploaderID   string
	Title        string
	Description  string
	Kind         string
	URL          string
	CreatedAt    time.Time
}
