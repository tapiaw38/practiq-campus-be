package assignment

import (
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/unlock"
)

type AssignmentData struct {
	ID              string  `json:"id"`
	CourseID        string  `json:"course_id"`
	SectionID       *string `json:"section_id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	DueAt           *string `json:"due_at"`
	MaxScore        int     `json:"max_score"`
	CreatedAt       string  `json:"created_at"`
	Weight          int     `json:"weight"`
	VisibleGroupID  *string `json:"visible_group_id"`
	UnlockAfterType *string `json:"unlock_after_type"`
	UnlockAfterID   *string `json:"unlock_after_id"`
	Locked          bool    `json:"locked"`
	LockedReason    string  `json:"locked_reason,omitempty"`
}

func formatDateTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02T15:04:05Z")
	return &s
}

func toAssignmentData(a domain.Assignment, status unlock.Status) AssignmentData {
	return AssignmentData{
		ID:              a.ID,
		CourseID:        a.CourseID,
		SectionID:       a.SectionID,
		Title:           a.Title,
		Description:     a.Description,
		DueAt:           formatDateTime(a.DueAt),
		MaxScore:        a.MaxScore,
		CreatedAt:       a.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Weight:          a.Weight,
		VisibleGroupID:  a.VisibleGroupID,
		UnlockAfterType: a.UnlockAfterType,
		UnlockAfterID:   a.UnlockAfterID,
		Locked:          status.Locked,
		LockedReason:    status.Reason,
	}
}
