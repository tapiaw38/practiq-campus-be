package domain

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	QuizQuestionMultipleChoice = "multiple_choice"
	QuizQuestionTrueFalse      = "true_false"
	QuizQuestionFillBlanks     = "fill_blanks"
)

type Quiz struct {
	ID             string
	CourseID       string
	SectionID      *string
	Title          string
	Description    string
	TimeLimitSecs  *int
	MaxAttempts    int
	ScheduledAt    *time.Time
	AvailableUntil *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time

	QuestionCount int

	Weight          int
	VisibleGroupID  *string
	UnlockAfterType *string
	UnlockAfterID   *string
}

type QuizQuestion struct {
	ID            string
	QuizID        string
	Type          string
	Statement     string
	Options       []string
	CorrectAnswer string
	Points        int
	Position      int
}

type QuizAttempt struct {
	ID            string
	QuizID        string
	UserID        string
	UserName      string
	AttemptNumber int
	StartedAt     time.Time

	ExpiresAt   *time.Time
	SubmittedAt *time.Time
	Score       int
	MaxScore    int
}

func (q Quiz) AcceptsNewAttempts(now time.Time) bool {
	if q.ScheduledAt != nil && now.Before(*q.ScheduledAt) {
		return false
	}
	if q.AvailableUntil != nil && now.After(*q.AvailableUntil) {
		return false
	}
	return true
}

func (a QuizAttempt) Expired(now time.Time) bool {
	return a.ExpiresAt != nil && now.After(*a.ExpiresAt)
}

func (a QuizAttempt) Live(now time.Time) bool {
	return a.SubmittedAt == nil && !a.Expired(now)
}

type QuizAnswer struct {
	ID         string
	AttemptID  string
	QuestionID string
	AnswerText string
	IsCorrect  bool
}

var blankMarkerPattern = regexp.MustCompile(`\{\{\s*(\d+)\s*\}\}`)

func BlankIDsInStatement(statement string) []int {
	matches := blankMarkerPattern.FindAllStringSubmatch(statement, -1)
	ids := make([]int, 0, len(matches))
	seen := map[int]bool{}
	for _, match := range matches {
		id, err := strconv.Atoi(match[1])
		if err != nil || seen[id] {
			continue
		}
		seen[id] = true
		ids = append(ids, id)
	}
	return ids
}

func NormalizeBlankAnswer(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func ParseBlanksAnswer(value string) map[string]string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	raw := map[string]string{}
	if err := json.Unmarshal([]byte(value), &raw); err != nil {
		return nil
	}
	parsed := make(map[string]string, len(raw))
	for id, answer := range raw {
		key := strings.TrimSpace(id)
		if number, err := strconv.Atoi(key); err == nil {
			key = strconv.Itoa(number)
		}
		parsed[key] = NormalizeBlankAnswer(answer)
	}
	return parsed
}

func BlanksAnswersMatch(student, expected string) bool {
	got := ParseBlanksAnswer(student)
	want := ParseBlanksAnswer(expected)
	if got == nil || want == nil || len(got) != len(want) {
		return false
	}
	for id, answer := range want {
		if !strings.EqualFold(got[id], answer) {
			return false
		}
	}
	return true
}

func ValidateQuizQuestion(q QuizQuestion) error {
	if strings.TrimSpace(q.Statement) == "" {
		return errInvalidQuestion("statement is required")
	}
	if q.Points <= 0 {
		return errInvalidQuestion("points must be positive")
	}
	switch q.Type {
	case QuizQuestionMultipleChoice:
		if len(q.Options) < 2 {
			return errInvalidQuestion("multiple choice needs at least 2 options")
		}
		found := false
		for _, opt := range q.Options {
			if strings.EqualFold(strings.TrimSpace(opt), strings.TrimSpace(q.CorrectAnswer)) {
				found = true
				break
			}
		}
		if !found {
			return errInvalidQuestion("correct_answer must be one of the options")
		}
	case QuizQuestionTrueFalse:
		answer := strings.ToLower(strings.TrimSpace(q.CorrectAnswer))
		if answer != "true" && answer != "false" {
			return errInvalidQuestion("correct_answer must be true or false")
		}
	case QuizQuestionFillBlanks:
		ids := BlankIDsInStatement(q.Statement)
		if len(ids) == 0 {
			return errInvalidQuestion("fill_blanks statement needs at least one {{n}} marker")
		}
		answers := ParseBlanksAnswer(q.CorrectAnswer)
		if len(answers) != len(ids) {
			return errInvalidQuestion("correct_answer must provide one answer per blank")
		}
		for _, id := range ids {
			answer, ok := answers[strconv.Itoa(id)]
			if !ok || answer == "" {
				return errInvalidQuestion("correct_answer is missing blank " + strconv.Itoa(id))
			}
		}
	default:
		return errInvalidQuestion("unsupported question type")
	}
	return nil
}

func GradeQuizAnswer(q QuizQuestion, answerText string) bool {
	switch q.Type {
	case QuizQuestionFillBlanks:
		return BlanksAnswersMatch(answerText, q.CorrectAnswer)
	default:
		return strings.EqualFold(strings.TrimSpace(answerText), strings.TrimSpace(q.CorrectAnswer))
	}
}

type invalidQuestionError string

func (e invalidQuestionError) Error() string { return string(e) }
func errInvalidQuestion(msg string) error    { return invalidQuestionError(msg) }
