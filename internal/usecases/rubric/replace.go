package rubric

import (
	"context"
	"strings"

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
	if e = a.Repositories.Rubric.Replace(c, id, out); e != nil {
		return apperrors.NewInternalError(e)
	}
	// The rubric's total *is* the maximum, so editing it moves the maximum
	// with it. Demanding they already match instead left no way to change
	// either: the rubric required the old maximum, and the maximum required
	// the old rubric total. Safe to do here because grading this assignment
	// has not started — the loop above refuses otherwise.
	if sum != task.MaxScore {
		task.MaxScore = sum
		if e = a.Repositories.Assignment.Update(c, id, *task); e != nil {
			return apperrors.NewInternalError(e)
		}
	}
	return nil
}
