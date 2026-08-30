package course_material

import (
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
)

// MaterialsFolder is the bucket prefix every course file lands under. It is
// also what the ownership check compares against, so it has to match between
// the upload call and the create call.
const MaterialsFolder = "materials"

const viewLinkTTL = time.Hour

type MaterialData struct {
	ID          string  `json:"id"`
	CourseID    string  `json:"course_id"`
	SectionID   *string `json:"section_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Kind        string  `json:"kind"`
	// URL is the stored value: a private bucket URL for files, or the
	// external link for kind="link".
	URL string `json:"url"`
	// ViewURL is what the browser actually opens — a short-lived signed URL
	// for files, and simply the same link for kind="link". Never stored.
	ViewURL   string `json:"view_url"`
	CreatedAt string `json:"created_at"`
}

func toMaterialData(m domain.CourseMaterial) MaterialData {
	return MaterialData{
		ID:          m.ID,
		CourseID:    m.CourseID,
		SectionID:   m.SectionID,
		Title:       m.Title,
		Description: m.Description,
		Kind:        m.Kind,
		URL:         m.URL,
		CreatedAt:   m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// withViewURL presigns a stored file so a reader can open it without the
// bucket being public. External links pass through untouched.
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
