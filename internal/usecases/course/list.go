package course

import (
	"context"

	courseRepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(context.Context, string, bool) (*ListOutput, apperrors.ApplicationError)
		ExecutePublished(context.Context) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []CourseData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) ExecutePublished(ctx context.Context) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	courses, err := app.Repositories.Course.List(ctx, courseRepo.ListFilter{PublishedOnly: true})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseListError, err)
	}
	return &ListOutput{Data: toCourseDataList(courses)}, nil
}

// Execute lists what the requester should see: a superadmin sees every
// course; a teacher sees the ones they own; a student sees the ones they are
// enrolled in. This is not "published courses browsing" (no course catalog
// exists yet in Phase 1) — it answers "my courses" for whoever asks.
func (u *listUsecase) Execute(ctx context.Context, requesterID string, isTeacher bool) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	// Campus profile type is source of truth for this app. A shared auth
	// account can carry the `admin` role for another product while still
	// being a student in Campus.
	if profile, err := app.Repositories.Profile.Get(ctx, requesterID); err == nil && profile != nil {
		isTeacher = profile.ProfileType == "teacher"
	}

	filter := courseRepo.ListFilter{}
	if isTeacher {
		filter.OwnerID = requesterID
	} else {
		filter.EnrolledUserID = requesterID
	}

	courses, err := app.Repositories.Course.List(ctx, filter)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseListError, err)
	}

	return &ListOutput{Data: toCourseDataList(courses)}, nil
}
