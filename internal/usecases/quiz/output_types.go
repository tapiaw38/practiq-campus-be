package quiz

import (
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type QuizData struct {
	ID             string  `json:"id"`
	CourseID       string  `json:"course_id"`
	SectionID      *string `json:"section_id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	TimeLimitSecs  *int    `json:"time_limit_secs"`
	MaxAttempts    int     `json:"max_attempts"`
	ScheduledAt    *string `json:"scheduled_at"`
	AvailableUntil *string `json:"available_until"`
	CreatedAt      string  `json:"created_at"`
}

func formatDateTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02T15:04:05Z")
	return &s
}

func toQuizData(q domain.Quiz) QuizData {
	return QuizData{
		ID:             q.ID,
		CourseID:       q.CourseID,
		SectionID:      q.SectionID,
		Title:          q.Title,
		Description:    q.Description,
		TimeLimitSecs:  q.TimeLimitSecs,
		MaxAttempts:    q.MaxAttempts,
		ScheduledAt:    formatDateTime(q.ScheduledAt),
		AvailableUntil: formatDateTime(q.AvailableUntil),
		CreatedAt:      q.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toQuizDataList(quizzes []domain.Quiz) []QuizData {
	out := make([]QuizData, 0, len(quizzes))
	for _, q := range quizzes {
		out = append(out, toQuizData(q))
	}
	return out
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
