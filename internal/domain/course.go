package domain

import "time"

const (
	CourseStatusDraft     = "draft"
	CourseStatusPublished = "published"
	CourseStatusArchived  = "archived"
)

type Course struct {
	ID          string
	TenantID    string
	OwnerID     string
	Title       string
	Slug        string
	Description string
	Status      string
	StartDate   *time.Time
	EndDate     *time.Time

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
