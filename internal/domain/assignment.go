package domain

import "time"

const (
	UnlockAfterAssignment = "assignment"
	UnlockAfterQuiz       = "quiz"
)

type Assignment struct {
	ID          string
	CourseID    string
	SectionID   *string
	Title       string
	Description string
	DueAt       *time.Time
	MaxScore    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	// Weight is this assignment's share of the course's final grade,
	// relative to every other weighted assignment/quiz — not a percentage
	// on its own. Two items both weighted 100 count equally regardless of
	// their max_score.
	Weight int
	// VisibleGroupID restricts this assignment to one course_groups cohort;
	// nil means every enrolled student sees it.
	VisibleGroupID *string
	// UnlockAfterType/UnlockAfterID name a single prerequisite item (an
	// assignment or a quiz) in the same course. Both nil means unlocked from
	// the start. A student meets it by having submitted the assignment or
	// submitted a quiz attempt — grading is not required, so a slow teacher
	// never permanently blocks a student.
	UnlockAfterType *string
	UnlockAfterID   *string
}
