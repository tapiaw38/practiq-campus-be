package assignment

import (
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type AssignmentData struct {
	ID          string  `json:"id"`
	CourseID    string  `json:"course_id"`
	SectionID   *string `json:"section_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	DueAt       *string `json:"due_at"`
	MaxScore    int     `json:"max_score"`
	CreatedAt   string  `json:"created_at"`
}

func formatDateTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02T15:04:05Z")
	return &s
}

func toAssignmentData(a domain.Assignment) AssignmentData {
	return AssignmentData{
		ID:          a.ID,
		CourseID:    a.CourseID,
		SectionID:   a.SectionID,
		Title:       a.Title,
		Description: a.Description,
		DueAt:       formatDateTime(a.DueAt),
		MaxScore:    a.MaxScore,
		CreatedAt:   a.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toAssignmentDataList(assignments []domain.Assignment) []AssignmentData {
	out := make([]AssignmentData, 0, len(assignments))
	for _, a := range assignments {
		out = append(out, toAssignmentData(a))
	}
	return out
}
