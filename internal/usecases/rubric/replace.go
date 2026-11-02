package rubric

import (
	"context"
	"errors"
	"strings"

	rubricrepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/rubric"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/usecases/campusaccess"
)

func (u *usecase) Replace(c context.Context, uid, id string, sa bool, in []Criterion) apperrors.ApplicationError {
	a := u.f()
	task, e := a.Repositories.Assignment.Get(c, id)
	if e != nil || task == nil {
		return apperrors.NewBadRequestError("assignment not found")
	}
	course, e := a.Repositories.Course.Get(c, task.CourseID)
	if e != nil {
		return apperrors.NewInternalError(e)
	}
	if course == nil || !campusaccess.CanManageCourse(c, course.OwnerID, uid, sa) {
		return apperrors.NewForbiddenError()
	}
	sum := 0
	out := make([]domain.RubricCriterion, 0, len(in))
	for _, v := range in {
		if strings.TrimSpace(v.Title) == "" || v.MaxScore < 1 {
			return apperrors.NewBadRequestError("invalid rubric criterion")
		}
		sum += v.MaxScore
		out = append(out, domain.RubricCriterion{Title: v.Title, Description: v.Description, MaxScore: v.MaxScore})
	}
	if len(out) == 0 {
		return apperrors.NewBadRequestError("la rúbrica necesita al menos un criterio")
	}

	submissions, e := a.Repositories.Submission.ListByAssignment(c, id)
	if e != nil {
		return apperrors.NewInternalError(e)
	}
	for _, submission := range submissions {
		if submission.Score != nil {
			return apperrors.NewBadRequestError("rubric cannot be changed after a submission is graded")
		}
	}
	if e = a.Repositories.Rubric.Replace(c, id, out, sum); e != nil {
		if errors.Is(e, rubricrepo.ErrAssignmentGraded) {
			return apperrors.NewBadRequestError("rubric cannot be changed after a submission is graded")
		}
		return apperrors.NewInternalError(e)
	}
	return nil
}
