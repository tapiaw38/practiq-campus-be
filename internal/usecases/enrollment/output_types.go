package enrollment

import "github.com/tapiaw38/practiq-campus-be/internal/domain"

type EnrollmentData struct {
	ID             string `json:"id"`
	CourseID       string `json:"course_id"`
	UserID         string `json:"user_id"`
	EnrollmentRole string `json:"enrollment_role"`
	Status         string `json:"status"`
	EnrolledAt     string `json:"enrolled_at"`
}

func toEnrollmentData(e domain.Enrollment) EnrollmentData {
	return EnrollmentData{
		ID:             e.ID,
		CourseID:       e.CourseID,
		UserID:         e.UserID,
		EnrollmentRole: e.EnrollmentRole,
		Status:         e.Status,
		EnrolledAt:     e.EnrolledAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toEnrollmentDataList(enrollments []domain.Enrollment) []EnrollmentData {
	out := make([]EnrollmentData, 0, len(enrollments))
	for _, e := range enrollments {
		out = append(out, toEnrollmentData(e))
	}
	return out
}
