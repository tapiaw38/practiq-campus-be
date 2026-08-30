package assignment

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
	x, e := a.Repositories.Assignment.Get(ctx, id)
	if e != nil {
		return apperrors.NewApplicationError(mappings.AssignmentGetError, e)
	}
	if x == nil {
		return apperrors.NewApplicationError(mappings.AssignmentNotFoundError, nil)
	}
	if z := requesterOwnsCourse(ctx, a, uid, sa, x.CourseID); z != nil {
		return z
	}
	if e = a.Repositories.Assignment.Delete(ctx, id); e != nil {
		return apperrors.NewApplicationError(mappings.AssignmentGetError, e)
	}
	return nil
}
