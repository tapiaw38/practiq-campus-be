package quiz

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
)

// validateUnlockTarget rejects a prerequisite that doesn't resolve: pointing
// at nothing, at itself, or at an item from a different course. Enforced
// here rather than only in the UI, since the API is the only real gate.
func validateUnlockTarget(ctx context.Context, app *appcontext.Context, courseID, selfID string, unlockType, unlockID *string) apperrors.ApplicationError {
	if unlockType == nil && unlockID == nil {
		return nil
	}
	if unlockType == nil || unlockID == nil || *unlockID == "" {
		return apperrors.NewBadRequestError("unlock_after requires both a type and an id")
	}
	if *unlockID == selfID {
		return apperrors.NewBadRequestError("an item cannot unlock after itself")
	}
	switch *unlockType {
	case domain.UnlockAfterAssignment:
		a, err := app.Repositories.Assignment.Get(ctx, *unlockID)
		if err != nil {
			return apperrors.NewInternalError(err)
		}
		if a == nil || a.CourseID != courseID {
			return apperrors.NewBadRequestError("unlock_after must reference an item in the same course")
		}
	case domain.UnlockAfterQuiz:
		q, err := app.Repositories.Quiz.Get(ctx, *unlockID)
		if err != nil {
			return apperrors.NewInternalError(err)
		}
		if q == nil || q.CourseID != courseID {
			return apperrors.NewBadRequestError("unlock_after must reference an item in the same course")
		}
	default:
		return apperrors.NewBadRequestError("unlock_after_type must be assignment or quiz")
	}
	return nil
}

func validateVisibleGroup(ctx context.Context, app *appcontext.Context, courseID string, groupID *string) apperrors.ApplicationError {
	if groupID == nil || *groupID == "" {
		return nil
	}
	g, err := app.Repositories.CourseGroup.Get(ctx, *groupID)
	if err != nil {
		return apperrors.NewInternalError(err)
	}
	if g == nil || g.CourseID != courseID {
		return apperrors.NewBadRequestError("visible_group_id must reference a group in the same course")
	}
	return nil
}
