package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	ListByCourseUsecase interface {
		Execute(context.Context, string, bool, string) (*ListOutput, apperrors.ApplicationError)
	}
	ListMineUsecase interface {
		Execute(context.Context, string) (*ListOutput, apperrors.ApplicationError)
	}

	listByCourseUsecase struct {
		contextFactory appcontext.Factory
	}
	listMineUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []EnrollmentData `json:"data"`
	}
)

func NewListByCourseUsecase(contextFactory appcontext.Factory) ListByCourseUsecase {
	return &listByCourseUsecase{contextFactory: contextFactory}
}

func NewListMineUsecase(contextFactory appcontext.Factory) ListMineUsecase {
	return &listMineUsecase{contextFactory: contextFactory}
}

func (u *listByCourseUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if appErr := requesterOwnsCourse(ctx, app, requesterID, isSuperAdmin, courseID); appErr != nil {
		return nil, appErr
	}

	enrollments, err := app.Repositories.Enrollment.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentListError, err)
	}

	return &ListOutput{Data: toEnrollmentDataList(enrollments)}, nil
}

// Execute answers the student dashboard: "which courses am I enrolled in."
func (u *listMineUsecase) Execute(ctx context.Context, requesterID string) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	enrollments, err := app.Repositories.Enrollment.ListByUser(ctx, requesterID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentListError, err)
	}

	return &ListOutput{Data: toEnrollmentDataList(enrollments)}, nil
}
