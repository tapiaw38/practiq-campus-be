package forum_post

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

// requesterCanAccessThread loads the thread and confirms the requester can
// see the course it belongs to (owner, superadmin, or enrolled) — a post
// can't be created/listed without first resolving which course gates it.
func requesterCanAccessThread(ctx context.Context, app *appcontext.Context, requesterID string, isSuperAdmin bool, threadID string) (*domain.ForumThread, apperrors.ApplicationError) {
	thread, err := app.Repositories.ForumThread.Get(ctx, threadID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ThreadGetError, err)
	}
	if thread == nil {
		return nil, apperrors.NewApplicationError(mappings.ThreadNotFoundError, nil)
	}

	if isSuperAdmin {
		return thread, nil
	}

	c, err := app.Repositories.Course.Get(ctx, thread.CourseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if c == nil {
		return nil, apperrors.NewApplicationError(mappings.CourseNotFoundError, nil)
	}
	if c.OwnerID == requesterID {
		return thread, nil
	}

	enrollment, err := app.Repositories.Enrollment.GetByCourseAndUser(ctx, thread.CourseID, requesterID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentGetError, err)
	}
	if enrollment == nil {
		return nil, apperrors.NewForbiddenError()
	}
	return thread, nil
}
