package notification

import (
	"context"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
)

type Usecase interface {
	List(context.Context, string) ([]domain.Notification, apperrors.ApplicationError)
	UnreadCount(context.Context, string) (int, apperrors.ApplicationError)
	Read(context.Context, string, string) apperrors.ApplicationError
}

func (u *usecase) UnreadCount(c context.Context, id string) (int, apperrors.ApplicationError) {
	count, err := u.f().Repositories.Notification.CountUnread(c, id)
	if err != nil {
		return 0, apperrors.NewInternalError(err)
	}
	return count, nil
}

type usecase struct{ f appcontext.Factory }

func New(f appcontext.Factory) Usecase { return &usecase{f} }
func (u *usecase) List(c context.Context, id string) ([]domain.Notification, apperrors.ApplicationError) {
	x, e := u.f().Repositories.Notification.List(c, id)
	if e != nil {
		return nil, apperrors.NewInternalError(e)
	}
	return x, nil
}
func (u *usecase) Read(c context.Context, id, user string) apperrors.ApplicationError {
	if e := u.f().Repositories.Notification.MarkRead(c, id, user); e != nil {
		return apperrors.NewInternalError(e)
	}
	return nil
}
