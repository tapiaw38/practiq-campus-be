package course_group

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type GroupData struct {
	ID        string   `json:"id"`
	CourseID  string   `json:"course_id"`
	Name      string   `json:"name"`
	CreatedAt string   `json:"created_at"`
	MemberIDs []string `json:"member_ids"`
}

func toGroupData(g domain.CourseGroup) GroupData {
	memberIDs := g.MemberIDs
	if memberIDs == nil {
		memberIDs = []string{}
	}
	return GroupData{ID: g.ID, CourseID: g.CourseID, Name: g.Name, CreatedAt: g.CreatedAt.Format("2006-01-02T15:04:05Z"), MemberIDs: memberIDs}
}

func toGroupDataList(groups []domain.CourseGroup) []GroupData {
	out := make([]GroupData, 0, len(groups))
	for _, g := range groups {
		out = append(out, toGroupData(g))
	}
	return out
}

type Usecase interface {
	List(ctx context.Context, courseID string) ([]GroupData, apperrors.ApplicationError)
	Create(ctx context.Context, requesterID string, isSuperAdmin bool, courseID, name string) (*GroupData, apperrors.ApplicationError)
	Delete(ctx context.Context, requesterID string, isSuperAdmin bool, groupID string) apperrors.ApplicationError
	AddMember(ctx context.Context, requesterID string, isSuperAdmin bool, groupID, userID string) apperrors.ApplicationError
	RemoveMember(ctx context.Context, requesterID string, isSuperAdmin bool, groupID, userID string) apperrors.ApplicationError
}

type usecase struct{ f appcontext.Factory }

func New(f appcontext.Factory) Usecase { return &usecase{f} }

func (u *usecase) requesterOwnsGroupsCourse(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string) apperrors.ApplicationError {
	if isSuperAdmin {
		return nil
	}
	app := u.f()
	c, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if c == nil {
		return apperrors.NewApplicationError(mappings.CourseNotFoundError, nil)
	}
	if c.OwnerID != requesterID {
		return apperrors.NewForbiddenError()
	}
	return nil
}

func (u *usecase) groupCourseID(ctx context.Context, groupID string) (string, apperrors.ApplicationError) {
	g, err := u.f().Repositories.CourseGroup.Get(ctx, groupID)
	if err != nil {
		return "", apperrors.NewInternalError(err)
	}
	if g == nil {
		return "", apperrors.NewBadRequestError("group not found")
	}
	return g.CourseID, nil
}
