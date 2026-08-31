package quiz

import (
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/unlock"
)

type QuizData struct {
	ID              string  `json:"id"`
	CourseID        string  `json:"course_id"`
	SectionID       *string `json:"section_id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	TimeLimitSecs   *int    `json:"time_limit_secs"`
	MaxAttempts     int     `json:"max_attempts"`
	ScheduledAt     *string `json:"scheduled_at"`
	AvailableUntil  *string `json:"available_until"`
	CreatedAt       string  `json:"created_at"`
	QuestionCount   int     `json:"question_count"`
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

func toQuizData(q domain.Quiz, status unlock.Status) QuizData {
	return QuizData{
		ID:              q.ID,
		CourseID:        q.CourseID,
		SectionID:       q.SectionID,
		Title:           q.Title,
		Description:     q.Description,
		TimeLimitSecs:   q.TimeLimitSecs,
		MaxAttempts:     q.MaxAttempts,
		ScheduledAt:     formatDateTime(q.ScheduledAt),
		AvailableUntil:  formatDateTime(q.AvailableUntil),
		CreatedAt:       q.CreatedAt.Format("2006-01-02T15:04:05Z"),
		QuestionCount:   q.QuestionCount,
		Weight:          q.Weight,
		VisibleGroupID:  q.VisibleGroupID,
		UnlockAfterType: q.UnlockAfterType,
		UnlockAfterID:   q.UnlockAfterID,
		Locked:          status.Locked,
		LockedReason:    status.Reason,
	}
}

// QuestionData is the teacher-facing shape: it carries CorrectAnswer, which
// must never reach a student before they submit an attempt.
type QuestionData struct {
	ID            string   `json:"id"`
	Type          string   `json:"type"`
	Statement     string   `json:"statement"`
	Options       []string `json:"options"`
	CorrectAnswer string   `json:"correct_answer"`
	Points        int      `json:"points"`
}

func toQuestionData(q domain.QuizQuestion) QuestionData {
	return QuestionData{
		ID:            q.ID,
		Type:          q.Type,
		Statement:     q.Statement,
		Options:       q.Options,
		CorrectAnswer: q.CorrectAnswer,
		Points:        q.Points,
	}
}

func toQuestionDataList(questions []domain.QuizQuestion) []QuestionData {
	out := make([]QuestionData, 0, len(questions))
	for _, q := range questions {
		out = append(out, toQuestionData(q))
	}
	return out
}
