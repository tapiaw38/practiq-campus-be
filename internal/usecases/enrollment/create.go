package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	CreateUsecase interface {
		Execute(context.Context, string, bool, string, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		UserID         string
		EnrollmentRole string
	}

	CreateOutput struct {
		Data EnrollmentData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	if input.UserID == "" {
		return nil, apperrors.NewBadRequestError("user_id is required")
	}

	app := u.contextFactory()

	if appErr := requesterOwnsCourse(ctx, app, requesterID, isSuperAdmin, courseID); appErr != nil {
		return nil, appErr
	}

	if existing, err := app.Repositories.Enrollment.GetByCourseAndUser(ctx, courseID, input.UserID); err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentGetError, err)
	} else if existing != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentAlreadyExistsError, nil)
	}

	role := input.EnrollmentRole
	if role == "" {
		role = domain.EnrollmentRoleStudent
	}

	id, err := app.Repositories.Enrollment.Create(ctx, domain.Enrollment{
		CourseID:       courseID,
		UserID:         input.UserID,
		EnrollmentRole: role,
		Status:         domain.EnrollmentStatusActive,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentCreateError, err)
	}

	created, err := app.Repositories.Enrollment.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentGetError, err)
	}
	if created == nil {
		return nil, apperrors.NewInternalError(nil)
	}

	return &CreateOutput{Data: toEnrollmentData(*created)}, nil
}
