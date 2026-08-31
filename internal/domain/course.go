package domain

import "time"

const (
	CourseStatusDraft     = "draft"
	CourseStatusPublished = "published"
	CourseStatusArchived  = "archived"
)

// Course is Campus's own course entity — deliberately unrelated to
// practiq-be's courses/topics/exercises. owner_id is the teacher of record
// for listing/audit; who else can act on the course beyond the owner is
// decided by Enrollment.Role, not by this struct.
type Course struct {
	ID          string
	OwnerID     string
	Title       string
	Slug        string
	Description string
	Status      string
	StartDate   *time.Time
	EndDate     *time.Time
	// PractiqSubjectID is an optional link to a practiq-be subject — a
	// reference only, never a copy: practiq's and Campus's course models
	// stay independent, this just lets a Campus course say "this
	// corresponds to that practiq subject" for reporting.
	PractiqSubjectID *string
	Labels           []string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CourseSection struct {
	ID          string
	CourseID    string
	Title       string
	Description string
	Position    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
