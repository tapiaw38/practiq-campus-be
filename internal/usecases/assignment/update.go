package assignment

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/unlock"
)

type UpdateUsecase interface {
	Execute(context.Context, string, bool, string, UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
}
type UpdateInput struct {
	SectionID          *string
	Title, Description string
	DueAt              *time.Time
	MaxScore           int
	Weight             int
	VisibleGroupID     *string
	UnlockAfterType    *string
	UnlockAfterID      *string
}
type UpdateOutput struct {
	Data AssignmentData `json:"data"`
}
type updateUsecase struct{ contextFactory appcontext.Factory }

func NewUpdateUsecase(f appcontext.Factory) UpdateUsecase { return &updateUsecase{f} }
func (u *updateUsecase) Execute(ctx context.Context, uid string, sa bool, id string, in UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	if strings.TrimSpace(in.Title) == "" {
		return nil, apperrors.NewBadRequestError("title is required")
	}
	a := u.contextFactory()
	x, e := a.Repositories.Assignment.Get(ctx, id)
	if e != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentGetError, e)
	}
	if x == nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentNotFoundError, nil)
	}
	if z := requesterOwnsCourse(ctx, a, uid, sa, x.CourseID); z != nil {
		return nil, z
	}
	if z := validateVisibleGroup(ctx, a, x.CourseID, in.VisibleGroupID); z != nil {
		return nil, z
	}
	if z := validateUnlockTarget(ctx, a, x.CourseID, id, in.UnlockAfterType, in.UnlockAfterID); z != nil {
		return nil, z
	}
	if in.MaxScore < 1 {
		return nil, apperrors.NewBadRequestError("el puntaje máximo debe ser mayor a cero")
	}

	if in.MaxScore != x.MaxScore {

		graded, e := a.Repositories.Submission.ListByAssignment(ctx, id)
		if e != nil {
			return nil, apperrors.NewInternalError(e)
		}
		for _, s := range graded {
			if s.GradedAt != nil {
				return nil, apperrors.NewBadRequestError(
					"ya hay entregas corregidas con este puntaje máximo: no se puede cambiar sin invalidar esas notas")
			}
		}

		criteria, e := a.Repositories.Rubric.List(ctx, id)
		if e != nil {
			return nil, apperrors.NewInternalError(e)
		}
		if len(criteria) > 0 {
			sum := 0
			for _, c := range criteria {
				sum += c.MaxScore
			}
			if sum != in.MaxScore {
				return nil, apperrors.NewBadRequestError(
					"la rúbrica define el puntaje máximo (" + strconv.Itoa(sum) + " puntos): cambialo editando la rúbrica")
			}
		}
	}

	weight := in.Weight
	if weight <= 0 {
		weight = 100
	}
	x.Title = in.Title
	x.Description = in.Description
	x.SectionID = in.SectionID
	x.DueAt = in.DueAt
	x.MaxScore = in.MaxScore
	x.Weight = weight
	x.VisibleGroupID = in.VisibleGroupID
	x.UnlockAfterType = in.UnlockAfterType
	x.UnlockAfterID = in.UnlockAfterID
	if e = a.Repositories.Assignment.Update(ctx, id, *x); e != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentGetError, e)
	}
	return &UpdateOutput{Data: toAssignmentData(*x, unlock.Status{})}, nil
}
