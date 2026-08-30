package calendar_event

import (
	"context"
	"sort"

	courseRepo "github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	ListMineUsecase interface {
		Execute(context.Context, string, bool) (*ListOutput, apperrors.ApplicationError)
	}

	listMineUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []EventData `json:"data"`
	}
)

func NewListMineUsecase(contextFactory appcontext.Factory) ListMineUsecase {
	return &listMineUsecase{contextFactory: contextFactory}
}

// Execute merges the requester's own manual events with the due dates of
// every assignment in a course they own or are enrolled in — the latter are
// assembled here, never stored, so there is nothing to keep in sync.
func (u *listMineUsecase) Execute(ctx context.Context, requesterID string, isTeacher bool) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	// Campus profile type wins over shared auth role (a student may be admin
	// in another Practiq product).
	if profile, err := app.Repositories.Profile.Get(ctx, requesterID); err == nil && profile != nil {
		isTeacher = profile.ProfileType == "teacher"
	}

	manual, err := app.Repositories.CalendarEvent.ListByOwner(ctx, requesterID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CalendarEventListError, err)
	}

	events := make([]EventData, 0, len(manual))
	for _, e := range manual {
		events = append(events, toEventData(e))
	}

	filter := courseRepo.ListFilter{}
	if isTeacher {
		filter.OwnerID = requesterID
	} else {
		filter.EnrolledUserID = requesterID
	}

	courses, err := app.Repositories.Course.List(ctx, filter)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseListError, err)
	}

	for _, c := range courses {
		assignments, err := app.Repositories.Assignment.ListByCourse(ctx, c.ID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.AssignmentListError, err)
		}
		for _, a := range assignments {
			if a.DueAt == nil {
				continue
			}
			courseID := c.ID
			events = append(events, EventData{
				ID:       "assignment:" + a.ID,
				CourseID: &courseID,
				Title:    "Vence: " + a.Title,
				StartsAt: formatTime(*a.DueAt),
				Source:   "assignment_due",
			})
		}
	}

	sort.Slice(events, func(i, j int) bool { return events[i].StartsAt < events[j].StartsAt })

	return &ListOutput{Data: events}, nil
}
