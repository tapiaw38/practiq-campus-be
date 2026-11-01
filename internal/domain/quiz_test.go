package domain

import (
	"testing"
	"time"
)

func TestGradeQuizAnswer(t *testing.T) {
	cases := []struct {
		name string
		q    QuizQuestion
		ans  string
		want bool
	}{
		{"multiple choice correct, case-insensitive", QuizQuestion{Type: QuizQuestionMultipleChoice, CorrectAnswer: "París"}, "parís", true},
		{"multiple choice wrong", QuizQuestion{Type: QuizQuestionMultipleChoice, CorrectAnswer: "París"}, "Roma", false},
		{"true_false correct", QuizQuestion{Type: QuizQuestionTrueFalse, CorrectAnswer: "true"}, "true", true},
		{"true_false wrong", QuizQuestion{Type: QuizQuestionTrueFalse, CorrectAnswer: "true"}, "false", false},
		{"fill_blanks correct despite spacing/case", QuizQuestion{Type: QuizQuestionFillBlanks, CorrectAnswer: `{"1":"estrella","2":"satélite"}`}, `{"1":"  Estrella ","2":"satélite"}`, true},
		{"fill_blanks missing a blank", QuizQuestion{Type: QuizQuestionFillBlanks, CorrectAnswer: `{"1":"estrella","2":"satélite"}`}, `{"1":"estrella"}`, false},
		{"fill_blanks malformed JSON", QuizQuestion{Type: QuizQuestionFillBlanks, CorrectAnswer: `{"1":"estrella"}`}, "not json", false},
		{"empty answer never matches", QuizQuestion{Type: QuizQuestionMultipleChoice, CorrectAnswer: "París"}, "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := GradeQuizAnswer(c.q, c.ans); got != c.want {
				t.Errorf("GradeQuizAnswer(%+v, %q) = %v, want %v", c.q, c.ans, got, c.want)
			}
		})
	}
}

func TestValidateQuizQuestion(t *testing.T) {
	cases := []struct {
		name    string
		q       QuizQuestion
		wantErr bool
	}{
		{"valid multiple choice", QuizQuestion{Type: QuizQuestionMultipleChoice, Statement: "s", Options: []string{"a", "b"}, CorrectAnswer: "a", Points: 1}, false},
		{"multiple choice correct_answer not in options", QuizQuestion{Type: QuizQuestionMultipleChoice, Statement: "s", Options: []string{"a", "b"}, CorrectAnswer: "c", Points: 1}, true},
		{"multiple choice single option", QuizQuestion{Type: QuizQuestionMultipleChoice, Statement: "s", Options: []string{"a"}, CorrectAnswer: "a", Points: 1}, true},
		{"valid true_false", QuizQuestion{Type: QuizQuestionTrueFalse, Statement: "s", CorrectAnswer: "true", Points: 1}, false},
		{"true_false bad answer", QuizQuestion{Type: QuizQuestionTrueFalse, Statement: "s", CorrectAnswer: "yes", Points: 1}, true},
		{"valid fill_blanks", QuizQuestion{Type: QuizQuestionFillBlanks, Statement: "a {{1}} b", CorrectAnswer: `{"1":"x"}`, Points: 1}, false},
		{"fill_blanks no markers", QuizQuestion{Type: QuizQuestionFillBlanks, Statement: "a b", CorrectAnswer: `{"1":"x"}`, Points: 1}, true},
		{"fill_blanks answer count mismatch", QuizQuestion{Type: QuizQuestionFillBlanks, Statement: "{{1}} {{2}}", CorrectAnswer: `{"1":"x"}`, Points: 1}, true},
		{"zero points", QuizQuestion{Type: QuizQuestionTrueFalse, Statement: "s", CorrectAnswer: "true", Points: 0}, true},
		{"empty statement", QuizQuestion{Type: QuizQuestionTrueFalse, Statement: "", CorrectAnswer: "true", Points: 1}, true},
		{"unsupported type", QuizQuestion{Type: "essay", Statement: "s", CorrectAnswer: "x", Points: 1}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateQuizQuestion(c.q)
			if (err != nil) != c.wantErr {
				t.Errorf("ValidateQuizQuestion(%+v) error = %v, wantErr %v", c.q, err, c.wantErr)
			}
		})
	}
}

func TestQuizAttemptLifecycle(t *testing.T) {
	now := time.Date(2026, 9, 17, 12, 0, 0, 0, time.UTC)
	past := now.Add(-time.Minute)
	future := now.Add(time.Minute)
	submitted := now.Add(-time.Hour)

	cases := []struct {
		name        string
		attempt     QuizAttempt
		wantExpired bool
		wantLive    bool
	}{
		{

			name:        "deadline passed",
			attempt:     QuizAttempt{ExpiresAt: &past},
			wantExpired: true,
		},
		{
			name:     "deadline still ahead",
			attempt:  QuizAttempt{ExpiresAt: &future},
			wantLive: true,
		},
		{

			name:     "no deadline is untimed, not expired",
			attempt:  QuizAttempt{},
			wantLive: true,
		},
		{
			name:    "already submitted is not live even before the deadline",
			attempt: QuizAttempt{ExpiresAt: &future, SubmittedAt: &submitted},
		},
		{
			name:        "exactly at the deadline is not yet over",
			attempt:     QuizAttempt{ExpiresAt: &now},
			wantExpired: false,
			wantLive:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.attempt.Expired(now); got != tc.wantExpired {
				t.Errorf("Expired = %v, want %v", got, tc.wantExpired)
			}
			if got := tc.attempt.Live(now); got != tc.wantLive {
				t.Errorf("Live = %v, want %v", got, tc.wantLive)
			}
		})
	}
}

func TestQuizAcceptsNewAttempts(t *testing.T) {
	now := time.Date(2026, 9, 17, 12, 0, 0, 0, time.UTC)
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	cases := []struct {
		name string
		quiz Quiz
		want bool
	}{
		{"no dates at all is always open", Quiz{}, true},
		{"opened in the past, no closing date", Quiz{ScheduledAt: &past}, true},
		{"not opened yet", Quiz{ScheduledAt: &future}, false},
		{"already closed", Quiz{AvailableUntil: &past}, false},
		{"inside the window", Quiz{ScheduledAt: &past, AvailableUntil: &future}, true},
		{"exactly at the opening instant counts as open", Quiz{ScheduledAt: &now}, true},
		{"exactly at the closing instant is not yet shut", Quiz{AvailableUntil: &now}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.quiz.AcceptsNewAttempts(now); got != tc.want {
				t.Fatalf("AcceptsNewAttempts = %v, want %v", got, tc.want)
			}
		})
	}
}
