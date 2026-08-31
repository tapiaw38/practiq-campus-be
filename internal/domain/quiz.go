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
	// QuestionCount is populated only by ListByCourse, for the card the
	// teacher sees before opening a quiz — it is never a stored column.
	QuestionCount int
	// Weight, VisibleGroupID, UnlockAfterType/UnlockAfterID mirror
	// Assignment's fields of the same name — see there for what each means.
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
	SubmittedAt   *time.Time
	Score         int
	MaxScore      int
}

type QuizAnswer struct {
	ID         string
	AttemptID  string
	QuestionID string
	AnswerText string
	IsCorrect  bool
}

// blankMarkerPattern matches {{n}} placeholders in a fill_blanks statement.
var blankMarkerPattern = regexp.MustCompile(`\{\{\s*(\d+)\s*\}\}`)

// BlankIDsInStatement returns the blank numbers declared in the statement.
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

// NormalizeBlankAnswer collapses whitespace so spacing never decides a grade.
func NormalizeBlankAnswer(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

// ParseBlanksAnswer reads the JSON answer format {"1":"a","2":"b"}.
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

// BlanksAnswersMatch reports whether the student placed every blank as
// expected. A missing, extra or malformed blank fails.
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

// ValidateQuizQuestion rejects a question the teacher could not intend, before
// it ever reaches a student. Enforced here (not just in the UI) since the API
// is the only gate.
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

// GradeQuizAnswer is the single source of truth for auto-grading: every
// supported question type is decided by exact comparison, deliberately never
// by AI, so a quiz result is reproducible and never waits on an external call.
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
