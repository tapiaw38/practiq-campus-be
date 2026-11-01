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

type Enrollment struct {
	ID             string
	CourseID       string
	UserID         string
	UserName       string
	EnrollmentRole string
	Status         string
	EnrolledAt     time.Time
}
