package rubric

import (
	"context"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"strings"
)

type Criterion struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MaxScore    int    `json:"max_score"`
}
type Usecase interface {
	List(context.Context, string, string, bool) ([]domain.RubricCriterion, apperrors.ApplicationError)
	Replace(context.Context, string, string, bool, []Criterion) apperrors.ApplicationError
}
type usecase struct{ f appcontext.Factory }

func New(f appcontext.Factory) Usecase { return &usecase{f} }
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
func (u *usecase) Replace(c context.Context, uid, id string, sa bool, in []Criterion) apperrors.ApplicationError {
	a := u.f()
	task, e := a.Repositories.Assignment.Get(c, id)
	if e != nil || task == nil {
		return apperrors.NewBadRequestError("assignment not found")
	}
	course, e := a.Repositories.Course.Get(c, task.CourseID)
	if e != nil || (!sa && course.OwnerID != uid) {
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
	if sum != task.MaxScore {
		return apperrors.NewBadRequestError("rubric total must equal max score")
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
	return nil
}
