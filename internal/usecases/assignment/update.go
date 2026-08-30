package assignment

import (
	"context"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"strings"
	"time"
)

type UpdateUsecase interface {
	Execute(context.Context, string, bool, string, UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
}
type UpdateInput struct {
	SectionID          *string
	Title, Description string
	DueAt              *time.Time
	MaxScore           int
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
	x.Title = in.Title
	x.Description = in.Description
	x.SectionID = in.SectionID
	x.DueAt = in.DueAt
	x.MaxScore = in.MaxScore
	if e = a.Repositories.Assignment.Update(ctx, id, *x); e != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentGetError, e)
	}
	return &UpdateOutput{Data: toAssignmentData(*x)}, nil
}
