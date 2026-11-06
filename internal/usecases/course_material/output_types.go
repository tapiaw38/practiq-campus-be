package course_material

import (
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
)

const MaterialsFolder = "materials"

const viewLinkTTL = time.Hour

type MaterialData struct {
	ID           string  `json:"id"`
	CourseID     string  `json:"course_id"`
	AssignmentID *string `json:"assignment_id"`
	SectionID    *string `json:"section_id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Kind         string  `json:"kind"`

	URL string `json:"url"`

	ViewURL   string `json:"view_url"`
	CreatedAt string `json:"created_at"`
}

func toMaterialData(m domain.CourseMaterial) MaterialData {
	return MaterialData{
		ID:           m.ID,
		CourseID:     m.CourseID,
		AssignmentID: m.AssignmentID,
		SectionID:    m.SectionID,
		Title:        m.Title,
		Description:  m.Description,
		Kind:         m.Kind,
		URL:          m.URL,
		CreatedAt:    m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func withViewURL(app *appcontext.Context, data MaterialData) MaterialData {
	if data.Kind == domain.MaterialKindLink {
		data.ViewURL = data.URL
		return data
	}
	if data.URL == "" || app.Storage == nil {
		return data
	}
	if signed, ok := app.Storage.PresignGetURL(data.URL, viewLinkTTL); ok {
		data.ViewURL = signed
	}
	return data
}

func toMaterialDataList(app *appcontext.Context, materials []domain.CourseMaterial) []MaterialData {
	out := make([]MaterialData, 0, len(materials))
	for _, m := range materials {
		out = append(out, withViewURL(app, toMaterialData(m)))
	}
	return out
}
