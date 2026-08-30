package rubric

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
)

func (u *usecase) List(c context.Context, id, uid string, sa bool) ([]domain.RubricCriterion, apperrors.ApplicationError) {
	a := u.f()
	task, e := a.Repositories.Assignment.Get(c, id)
	if e != nil || task == nil {
		return nil, apperrors.NewBadRequestError("assignment not found")
	}
	if !sa {
		course, e := a.Repositories.Course.Get(c, task.CourseID)
		if e != nil {
			return nil, apperrors.NewInternalError(e)
		}
		if course == nil {
			return nil, apperrors.NewBadRequestError("course not found")
		}
		if course.OwnerID != uid {
			enrollment, e := a.Repositories.Enrollment.GetByCourseAndUser(c, task.CourseID, uid)
			if e != nil {
				return nil, apperrors.NewInternalError(e)
			}
			if enrollment == nil {
				return nil, apperrors.NewForbiddenError()
			}
		}
	}
	x, e := a.Repositories.Rubric.List(c, id)
	if e != nil {
		return nil, apperrors.NewInternalError(e)
	}
	return x, nil
}
