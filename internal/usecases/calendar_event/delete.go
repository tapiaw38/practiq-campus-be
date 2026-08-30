package calendar_event

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type DeleteUsecase interface {
	Execute(context.Context, string, string) apperrors.ApplicationError
}
type deleteUsecase struct{ contextFactory appcontext.Factory }

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{contextFactory: contextFactory}
}
func (u *deleteUsecase) Execute(ctx context.Context, requesterID, eventID string) apperrors.ApplicationError {
	app := u.contextFactory()
	event, err := app.Repositories.CalendarEvent.Get(ctx, eventID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.CalendarEventGetError, err)
	}
	if event == nil {
		return apperrors.NewNotFoundError("event not found")
	}
	if event.OwnerID != requesterID {
		return apperrors.NewForbiddenError()
	}
	if err := app.Repositories.CalendarEvent.Delete(ctx, eventID); err != nil {
		return apperrors.NewApplicationError(mappings.CalendarEventDeleteError, err)
	}
	for _, attendeeID := range event.AttendeeIDs {
		_ = app.Repositories.Notification.Create(ctx, domain.Notification{UserID: attendeeID, Type: "calendar_event_cancelled", Title: "Evento cancelado: " + event.Title, Body: "El evento fue eliminado", Data: `{"event_id":"` + eventID + `"}`})
	}
	return nil
}
