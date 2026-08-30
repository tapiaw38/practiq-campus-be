package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type SelfUsecase interface {
	Execute(context.Context, string, string) (*CreateOutput, apperrors.ApplicationError)
}
type selfUsecase struct{ contextFactory appcontext.Factory }

func NewSelfUsecase(f appcontext.Factory) SelfUsecase { return &selfUsecase{f} }
func (u *selfUsecase) Execute(ctx context.Context, uid, cid string) (*CreateOutput, apperrors.ApplicationError) {
	a := u.contextFactory()
	c, e := a.Repositories.Course.Get(ctx, cid)
	if e != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, e)
	}
	if c == nil {
		return nil, apperrors.NewApplicationError(mappings.CourseNotFoundError, nil)
	}
	if c.Status != domain.CourseStatusPublished {
		return nil, apperrors.NewBadRequestError("course is not published")
	}
	if x, e := a.Repositories.Enrollment.GetByCourseAndUser(ctx, cid, uid); e != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentGetError, e)
	} else if x != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentAlreadyExistsError, nil)
	}
	id, e := a.Repositories.Enrollment.Create(ctx, domain.Enrollment{CourseID: cid, UserID: uid, EnrollmentRole: domain.EnrollmentRoleStudent, Status: domain.EnrollmentStatusActive})
	if e != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentCreateError, e)
	}
	x, e := a.Repositories.Enrollment.Get(ctx, id)
	if e != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentGetError, e)
	}
	if x == nil {
		return nil, apperrors.NewInternalError(nil)
	}
	return &CreateOutput{Data: toEnrollmentData(*x)}, nil
}
