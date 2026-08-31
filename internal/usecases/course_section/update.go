package course_section

import (
	"context"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"strings"
)

type UpdateUsecase interface {
	Execute(context.Context, string, bool, string, UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
}
type UpdateInput struct{ Title, Description string }
type UpdateOutput struct {
	Data SectionData `json:"data"`
}
type updateUsecase struct{ contextFactory appcontext.Factory }

func NewUpdateUsecase(f appcontext.Factory) UpdateUsecase { return &updateUsecase{f} }
func (u *updateUsecase) Execute(ctx context.Context, uid string, sa bool, id string, in UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	if strings.TrimSpace(in.Title) == "" {
		return nil, apperrors.NewBadRequestError("title is required")
	}
	a := u.contextFactory()
	s, e := a.Repositories.CourseSection.Get(ctx, id)
	if e != nil {
		return nil, apperrors.NewApplicationError(mappings.SectionGetError, e)
	}
	if s == nil {
		return nil, apperrors.NewApplicationError(mappings.CourseNotFoundError, nil)
	}
	if x := requesterOwnsCourse(ctx, a, uid, sa, s.CourseID); x != nil {
		return nil, x
	}
	s.Title = in.Title
	s.Description = in.Description
	if e = a.Repositories.CourseSection.Update(ctx, id, *s); e != nil {
		return nil, apperrors.NewApplicationError(mappings.SectionGetError, e)
	}
	return &UpdateOutput{Data: toSectionData(*s)}, nil
}
