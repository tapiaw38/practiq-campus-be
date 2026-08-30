package quiz_attempt

import (
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func formatDateTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02T15:04:05Z")
	return &s
}

// StudentQuestionData is what a student sees while an attempt is in
// progress: never CorrectAnswer, or the quiz would grade itself.
type StudentQuestionData struct {
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	Statement string   `json:"statement"`
	Options   []string `json:"options"`
	Points    int      `json:"points"`
}

func toStudentQuestionData(questions []domain.QuizQuestion) []StudentQuestionData {
	out := make([]StudentQuestionData, 0, len(questions))
	for _, q := range questions {
		out = append(out, StudentQuestionData{ID: q.ID, Type: q.Type, Statement: q.Statement, Options: q.Options, Points: q.Points})
	}
	return out
}

type AttemptData struct {
	ID            string  `json:"id"`
	QuizID        string  `json:"quiz_id"`
	UserID        string  `json:"user_id"`
	UserName      string  `json:"user_name,omitempty"`
	AttemptNumber int     `json:"attempt_number"`
	StartedAt     string  `json:"started_at"`
	SubmittedAt   *string `json:"submitted_at"`
	Score         int     `json:"score"`
	MaxScore      int     `json:"max_score"`
}

func toAttemptData(a domain.QuizAttempt) AttemptData {
	return AttemptData{
		ID:            a.ID,
		QuizID:        a.QuizID,
		UserID:        a.UserID,
		UserName:      a.UserName,
		AttemptNumber: a.AttemptNumber,
		StartedAt:     a.StartedAt.Format("2006-01-02T15:04:05Z"),
		SubmittedAt:   formatDateTime(a.SubmittedAt),
		Score:         a.Score,
		MaxScore:      a.MaxScore,
	}
}

func toAttemptDataList(attempts []domain.QuizAttempt) []AttemptData {
	out := make([]AttemptData, 0, len(attempts))
	for _, a := range attempts {
		out = append(out, toAttemptData(a))
	}
	return out
}

// AnswerResultData reveals the correct answer: only ever returned once the
// attempt holding it has been submitted.
type AnswerResultData struct {
	QuestionID    string `json:"question_id"`
	Statement     string `json:"statement"`
	AnswerText    string `json:"answer_text"`
	CorrectAnswer string `json:"correct_answer"`
	IsCorrect     bool   `json:"is_correct"`
	Points        int    `json:"points"`
}
