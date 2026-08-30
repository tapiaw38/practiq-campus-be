package course_section

import (
	"context"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type DeleteUsecase interface {
	Execute(context.Context, string, bool, string) apperrors.ApplicationError
}
type deleteUsecase struct{ contextFactory appcontext.Factory }

func NewDeleteUsecase(f appcontext.Factory) DeleteUsecase { return &deleteUsecase{f} }
func (u *deleteUsecase) Execute(ctx context.Context, uid string, sa bool, id string) apperrors.ApplicationError {
	a := u.contextFactory()
	s, e := a.Repositories.CourseSection.Get(ctx, id)
	if e != nil {
		return apperrors.NewApplicationError(mappings.SectionGetError, e)
	}
	if s == nil {
		return apperrors.NewApplicationError(mappings.CourseNotFoundError, nil)
	}
	if x := requesterOwnsCourse(ctx, a, uid, sa, s.CourseID); x != nil {
		return x
	}
	if e = a.Repositories.CourseSection.Delete(ctx, id); e != nil {
		return apperrors.NewApplicationError(mappings.SectionGetError, e)
	}
	return nil
}
