// Package unlock resolves the single prerequisite an assignment or quiz may
// declare (UnlockAfterType/UnlockAfterID) against one student's progress.
// It is shared because both assignment and quiz usecases need the same
// answer — listing (to flag an item locked) and creating a submission or
// starting an attempt (to actually reject it), so the rule is defined once.
package unlock

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type Status struct {
	Locked bool
	Reason string
}

// Check meets the prerequisite with a submission (assignment) or a submitted
// attempt (quiz) — never a passing grade, so a slow teacher can never
// permanently lock a student out of the rest of the course.
func Check(ctx context.Context, repos *repositories.Repositories, unlockAfterType, unlockAfterID *string, studentID string) (Status, error) {
	if unlockAfterType == nil || unlockAfterID == nil || studentID == "" {
		return Status{}, nil
	}

	switch *unlockAfterType {
	case domain.UnlockAfterAssignment:
		a, err := repos.Assignment.Get(ctx, *unlockAfterID)
		if err != nil || a == nil {
			return Status{}, err
		}
		submission, err := repos.Submission.GetByAssignmentAndUser(ctx, *unlockAfterID, studentID)
		if err != nil {
			return Status{}, err
		}
		if submission != nil {
			return Status{}, nil
		}
		return Status{Locked: true, Reason: "Disponible después de entregar: " + a.Title}, nil

	case domain.UnlockAfterQuiz:
		q, err := repos.Quiz.Get(ctx, *unlockAfterID)
		if err != nil || q == nil {
			return Status{}, err
		}
		attempts, err := repos.QuizAttempt.ListMine(ctx, *unlockAfterID, studentID)
		if err != nil {
			return Status{}, err
		}
		for _, attempt := range attempts {
			if attempt.SubmittedAt != nil {
				return Status{}, nil
			}
		}
		return Status{Locked: true, Reason: "Disponible después de rendir: " + q.Title}, nil

	default:
		return Status{}, nil
	}
}
