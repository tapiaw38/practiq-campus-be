package submission

import "github.com/tapiaw38/practiq-campus-be/internal/domain"

type SubmissionData struct {
	ID           string  `json:"id"`
	AssignmentID string  `json:"assignment_id"`
	UserID       string  `json:"user_id"`
	UserName     string  `json:"user_name"`
	Content      string  `json:"content"`
	Status       string  `json:"status"`
	Score        *int    `json:"score"`
	Feedback     string  `json:"feedback"`
	SubmittedAt  string  `json:"submitted_at"`
	GradedAt     *string `json:"graded_at"`
}

func toSubmissionData(s domain.Submission) SubmissionData {
	var gradedAt *string
	if s.GradedAt != nil {
		v := s.GradedAt.Format("2006-01-02T15:04:05Z")
		gradedAt = &v
	}
	return SubmissionData{
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
