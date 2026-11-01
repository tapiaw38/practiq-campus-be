package domain

import "time"

const (
	MaterialKindFile = "file"
	MaterialKindLink = "link"
)

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
