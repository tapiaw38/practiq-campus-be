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

// Execute lists what the requester should see: a superadmin sees every
// course; a teacher sees the ones they own; a student sees the ones they are
// enrolled in. This is not "published courses browsing" (no course catalog
// exists yet in Phase 1) — it answers "my courses" for whoever asks.
func (u *listUsecase) Execute(ctx context.Context, requesterID string, isTeacher bool) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

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
