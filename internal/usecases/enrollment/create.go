package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/identity"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

type (
	CreateUsecase interface {
		Execute(context.Context, string, bool, string, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		Email          string
		EnrollmentRole string
		BearerToken    string
	}

	CreateOutput struct {
		Data EnrollmentData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	if input.Email == "" {
		return nil, apperrors.NewBadRequestError("email is required")
	}

	app := u.contextFactory()

	if appErr := requesterOwnsCourse(ctx, app, requesterID, isSuperAdmin, courseID); appErr != nil {
		return nil, appErr
	}

	authUser, err := app.Integrations.AuthAPI.GetByEmail(ctx, input.BearerToken, input.Email)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if authUser == nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentStudentNotFoundError, nil)
	}
	student, err := app.Repositories.Profile.Get(ctx, authUser.Username)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if student == nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentStudentNotFoundError, nil)
	}

	if appErr := studentBelongsToInstitution(ctx, app, input.BearerToken, student.ID); appErr != nil {
		return nil, appErr
	}

	if existing, err := app.Repositories.Enrollment.GetByCourseAndUser(ctx, courseID, student.ID); err != nil {
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
		UserID:         student.ID,
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

	names, err := identity.Names(ctx, app.Integrations.AuthAPI, input.BearerToken, []string{created.UserID})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	created.UserName = identity.FullName(names[created.UserID], created.UserID)

	return &CreateOutput{Data: toEnrollmentData(*created)}, nil
}

func studentBelongsToInstitution(ctx context.Context, app *appcontext.Context, bearerToken, studentID string) apperrors.ApplicationError {
	tenantID := tenantcontext.ID(ctx)
	if tenantID == "" {
		return apperrors.NewInternalError(tenantcontext.ErrNoTenant)
	}
	tenant, err := app.Repositories.Tenant.Get(ctx, tenantID)
	if err != nil {
		return apperrors.NewInternalError(err)
	}
	if tenant == nil {
		return apperrors.NewInternalError(nil)
	}

	members, err := app.Integrations.PractiqAPI.ListSchoolMembers(ctx, bearerToken, tenant.SchoolID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.PractiqLookupError, err)
	}
	for _, member := range members {
		if member.UserID == studentID && member.Active {
			return nil
		}
	}
	return apperrors.NewApplicationError(mappings.EnrollmentStudentNotInInstitutionError, nil)
}
