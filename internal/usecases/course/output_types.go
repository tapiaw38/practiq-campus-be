package course

import (
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type CourseData struct {
	ID          string  `json:"id"`
	OwnerID     string  `json:"owner_id"`
	Title       string  `json:"title"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func formatDate(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

func toCourseData(c domain.Course) CourseData {
	return CourseData{
		ID:          c.ID,
		OwnerID:     c.OwnerID,
		Title:       c.Title,
		Slug:        c.Slug,
		Description: c.Description,
		Status:      c.Status,
		StartDate:   formatDate(c.StartDate),
		EndDate:     formatDate(c.EndDate),
		CreatedAt:   c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toCourseDataList(courses []domain.Course) []CourseData {
	out := make([]CourseData, 0, len(courses))
	for _, c := range courses {
		out = append(out, toCourseData(c))
	}
	return out
}
