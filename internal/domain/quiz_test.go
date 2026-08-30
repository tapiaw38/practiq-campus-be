package domain

import "testing"

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
