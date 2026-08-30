package course_section

import "github.com/tapiaw38/practiq-campus-be/internal/domain"

type SectionData struct {
	ID        string `json:"id"`
	CourseID  string `json:"course_id"`
	Title     string `json:"title"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
}

func toSectionData(s domain.CourseSection) SectionData {
	return SectionData{
		ID:        s.ID,
		CourseID:  s.CourseID,
		Title:     s.Title,
		Position:  s.Position,
		CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toSectionDataList(sections []domain.CourseSection) []SectionData {
	out := make([]SectionData, 0, len(sections))
	for _, s := range sections {
		out = append(out, toSectionData(s))
	}
	return out
}
