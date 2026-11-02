package submission

import (
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
)

const SubmissionsFolder = "submissions"

const attachmentLinkTTL = time.Hour

const attachmentNameMarker = "Archivo adjunto:"

type AttachmentData struct {
	Filename string `json:"filename"`

	URL string `json:"url"`

	ViewURL string `json:"view_url"`
}

type RubricScoreData struct {
	CriterionID string `json:"criterion_id"`
	Score       int    `json:"score"`
	Feedback    string `json:"feedback"`
}

type SubmissionData struct {
	ID           string            `json:"id"`
	AssignmentID string            `json:"assignment_id"`
	UserID       string            `json:"user_id"`
	UserName     string            `json:"user_name"`
	Content      string            `json:"content"`
	Status       string            `json:"status"`
	Version      int               `json:"version"`
	Score        *int              `json:"score"`
	Feedback     string            `json:"feedback"`
	SubmittedAt  string            `json:"submitted_at"`
	GradedAt     *string           `json:"graded_at"`
	RubricScores []RubricScoreData `json:"rubric_scores"`
	Attachments  []AttachmentData  `json:"attachments"`
}

func withRubricScores(data SubmissionData, scores []domain.RubricScore) SubmissionData {
	data.RubricScores = make([]RubricScoreData, 0, len(scores))
	for _, score := range scores {
		data.RubricScores = append(data.RubricScores, RubricScoreData{CriterionID: score.CriterionID, Score: score.Score, Feedback: score.Feedback})
	}
	return data
}

func withAttachments(app *appcontext.Context, data SubmissionData) SubmissionData {
	if app == nil || app.Storage == nil || data.Content == "" {
		return data
	}

	pendingName := ""
	for _, line := range strings.Split(data.Content, "\n") {
		trimmed := strings.TrimSpace(line)
		if name, found := strings.CutPrefix(trimmed, attachmentNameMarker); found {
			pendingName = strings.TrimSpace(name)
			continue
		}
		for _, token := range strings.Fields(trimmed) {
			if !app.Storage.OwnsFileURL(token, SubmissionsFolder, data.UserID) {
				continue
			}
			signed, ok := app.Storage.PresignGetURL(token, attachmentLinkTTL)
			if !ok {
				continue
			}
			data.Attachments = append(data.Attachments, AttachmentData{
				Filename: attachmentFilename(pendingName, token),
				URL:      token,
				ViewURL:  signed,
			})
			pendingName = ""
		}
	}

	return data
}

func attachmentFilename(preferred, rawURL string) string {
	if preferred != "" {
		return preferred
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	name := path.Base(parsed.Path)
	if unescaped, err := url.PathUnescape(name); err == nil {
		name = unescaped
	}
	if name == "" || name == "." || name == "/" {
		return rawURL
	}
	return name
}

func toSubmissionData(s domain.Submission) SubmissionData {
	var gradedAt *string
	if s.GradedAt != nil {
		v := s.GradedAt.Format("2006-01-02T15:04:05Z")
		gradedAt = &v
	}
	return SubmissionData{
		Attachments:  make([]AttachmentData, 0),
		ID:           s.ID,
		AssignmentID: s.AssignmentID,
		UserID:       s.UserID,
		UserName:     s.UserName,
		Content:      s.Content,
		Status:       s.Status,
		Version:      s.Version,
		Score:        s.Score,
		Feedback:     s.Feedback,
		SubmittedAt:  s.SubmittedAt.Format("2006-01-02T15:04:05Z"),
		GradedAt:     gradedAt,
	}
}

func toSubmissionDataList(submissions []domain.Submission) []SubmissionData {
	out := make([]SubmissionData, 0, len(submissions))
	for _, s := range submissions {
		out = append(out, toSubmissionData(s))
	}
	return out
}
