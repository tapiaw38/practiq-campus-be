package submission

import (
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
)

// SubmissionsFolder is the bucket prefix a student's submission files land
// under. It has to match the folder the upload call used, because that is what
// the ownership check compares against.
const SubmissionsFolder = "submissions"

const attachmentLinkTTL = time.Hour

// attachmentNameMarker is the line the frontend writes above an uploaded
// file's URL. It only supplies a nicer display name — an attachment is
// recognised by its URL, so a submission that lost the marker still resolves.
const attachmentNameMarker = "Archivo adjunto:"

// AttachmentData is a file the student uploaded with their submission.
type AttachmentData struct {
	Filename string `json:"filename"`
	// URL is the stored private-bucket value, kept so the client can tell two
	// attachments apart across refreshes; it is not openable on its own.
	URL string `json:"url"`
	// ViewURL is what the browser opens — a short-lived signed URL. Never
	// stored.
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

// withAttachments presigns the files a submission references so the teacher
// grading it — and the student re-reading it — can actually open them. Files
// live in a private bucket, so the URL stored in the submission body is not
// openable; materials solve this the same way, in course_material.withViewURL.
//
// Content is free text the student wrote, so a URL appearing in it proves
// nothing: only objects the submitting student uploaded are signed. Signing
// every bucket URL found in the body would turn the submission form into a
// way to read anyone else's files.
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

// attachmentFilename prefers the name the student's browser reported and falls
// back to the object name, which carries an upload prefix but still beats
// showing a bare URL.
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
