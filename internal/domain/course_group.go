package domain

import "time"

// CourseGroup partitions a course's enrolled students into a named cohort
// (e.g. "Comisión A") for organizing and filtering, independent of
// CourseSection which orders content instead of people.
type CourseGroup struct {
	ID        string
	CourseID  string
	Name      string
	CreatedAt time.Time
	MemberIDs []string
}
