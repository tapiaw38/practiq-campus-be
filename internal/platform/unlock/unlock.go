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
