package domain

import "time"

const (
	EnrollmentRoleStudent    = "student"
	EnrollmentRoleTA         = "teaching_assistant"
	EnrollmentRoleCoTeacher  = "co_teacher"
	EnrollmentStatusActive   = "active"
	EnrollmentStatusDropped  = "dropped"
	EnrollmentStatusComplete = "completed"
)

// Enrollment joins a Profile to a Course. Role is course-scoped: it grants
// elevated access (teaching_assistant/co_teacher) within THIS course only,
// independent of the user's global JWT role and independent of who owns the
// course — this is how a student in one course can TA another without a new
// global role or a course-level RBAC table.
type Enrollment struct {
	ID             string
	CourseID       string
	UserID         string
	EnrollmentRole string
	Status         string
	EnrolledAt     time.Time
}
