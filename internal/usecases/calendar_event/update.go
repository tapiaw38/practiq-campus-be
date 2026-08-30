package calendar_event

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type UpdateUsecase interface {
	Execute(context.Context, string, string, CreateInput) (*UpdateOutput, apperrors.ApplicationError)
}
type UpdateOutput struct {
	Data EventData `json:"data"`
}
type updateUsecase struct{ contextFactory appcontext.Factory }

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, requesterID, eventID string, input CreateInput) (*UpdateOutput, apperrors.ApplicationError) {
	if appErr := validateInput(input); appErr != nil {
		return nil, appErr
	}
	app := u.contextFactory()
	existing, err := app.Repositories.CalendarEvent.Get(ctx, eventID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CalendarEventGetError, err)
	}
	if existing == nil {
		return nil, apperrors.NewNotFoundError("event not found")
	}
	if existing.OwnerID != requesterID {
		return nil, apperrors.NewForbiddenError()
	}
	attendees, appErr := validateCourseAndAttendees(ctx, app, requesterID, input.CourseID, input.AttendeeIDs)
	if appErr != nil {
		return nil, appErr
	}
	if err := app.Repositories.CalendarEvent.Update(ctx, eventID, domain.CalendarEvent{
		CourseID: input.CourseID, Title: strings.TrimSpace(input.Title), Description: strings.TrimSpace(input.Description), StartsAt: input.StartsAt, EndsAt: input.EndsAt, AllDay: input.AllDay, RecurrenceRule: normalizeRecurrence(input.RecurrenceRule), ReminderMinutes: input.ReminderMinutes, AttendeeIDs: attendees,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.CalendarEventUpdateError, err)
	}
	updated, err := app.Repositories.CalendarEvent.Get(ctx, eventID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CalendarEventGetError, err)
	}
	if updated == nil {
		return nil, apperrors.NewInternalError(nil)
	}
	return &UpdateOutput{Data: toEventData(*updated)}, nil
}

func validateInput(input CreateInput) apperrors.ApplicationError {
	if strings.TrimSpace(input.Title) == "" {
		return apperrors.NewBadRequestError("title is required")
	}
	if input.CourseID == nil || strings.TrimSpace(*input.CourseID) == "" {
		return apperrors.NewBadRequestError("course_id is required")
	}
	if input.EndsAt != nil && !input.EndsAt.After(input.StartsAt) {
		return apperrors.NewBadRequestError("ends_at must be after starts_at")
	}
	rule := normalizeRecurrence(input.RecurrenceRule)
	if rule != "none" && rule != "daily" && rule != "weekly" && rule != "monthly" {
		return apperrors.NewBadRequestError("invalid recurrence_rule")
	}
	if input.ReminderMinutes != nil && (*input.ReminderMinutes < 0 || *input.ReminderMinutes > 10080) {
		return apperrors.NewBadRequestError("invalid reminder_minutes")
	}
	return nil
}

func normalizeRecurrence(rule string) string {
	if rule == "" {
		return "none"
	}
	return rule
}

func validateCourseAndAttendees(ctx context.Context, app *appcontext.Context, requesterID string, courseID *string, attendeeIDs []string) ([]string, apperrors.ApplicationError) {
	course, err := app.Repositories.Course.Get(ctx, *courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if course == nil {
		return nil, apperrors.NewNotFoundError("course not found")
	}
	if course.OwnerID != requesterID {
		return nil, apperrors.NewForbiddenError()
	}
	enrollments, err := app.Repositories.Enrollment.ListByCourse(ctx, *courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentListError, err)
	}
	allowed, seen := make(map[string]bool, len(enrollments)), make(map[string]bool)
	for _, enrollment := range enrollments {
		if enrollment.Status == domain.EnrollmentStatusActive {
			allowed[enrollment.UserID] = true
		}
	}
	attendees := make([]string, 0, len(attendeeIDs))
	for _, userID := range attendeeIDs {
		if !allowed[userID] {
			return nil, apperrors.NewBadRequestError("attendee must be enrolled in the selected course")
		}
		if !seen[userID] {
			attendees = append(attendees, userID)
			seen[userID] = true
		}
	}
	return attendees, nil
}
