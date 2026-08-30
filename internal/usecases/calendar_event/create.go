package calendar_event

import (
	"context"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	CreateUsecase interface {
		Execute(context.Context, string, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		CourseID        *string
		AttendeeIDs     []string
		Title           string
		Description     string
		StartsAt        time.Time
		EndsAt          *time.Time
		AllDay          bool
		RecurrenceRule  string
		ReminderMinutes *int
	}

	CreateOutput struct {
		Data EventData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, requesterID string, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, apperrors.NewBadRequestError("title is required")
	}
	if input.CourseID == nil || strings.TrimSpace(*input.CourseID) == "" {
		return nil, apperrors.NewBadRequestError("course_id is required")
	}
	if input.EndsAt != nil && !input.EndsAt.After(input.StartsAt) {
		return nil, apperrors.NewBadRequestError("ends_at must be after starts_at")
	}
	if input.RecurrenceRule == "" {
		input.RecurrenceRule = "none"
	}
	if input.RecurrenceRule != "none" && input.RecurrenceRule != "daily" && input.RecurrenceRule != "weekly" && input.RecurrenceRule != "monthly" {
		return nil, apperrors.NewBadRequestError("invalid recurrence_rule")
	}
	if input.ReminderMinutes != nil && (*input.ReminderMinutes < 0 || *input.ReminderMinutes > 10080) {
		return nil, apperrors.NewBadRequestError("invalid reminder_minutes")
	}

	app := u.contextFactory()
	course, err := app.Repositories.Course.Get(ctx, *input.CourseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if course == nil {
		return nil, apperrors.NewNotFoundError("course not found")
	}
	if course.OwnerID != requesterID {
		return nil, apperrors.NewForbiddenError()
	}
	enrollments, err := app.Repositories.Enrollment.ListByCourse(ctx, *input.CourseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentListError, err)
	}
	allowed := make(map[string]bool, len(enrollments))
	for _, enrollment := range enrollments {
		if enrollment.Status == domain.EnrollmentStatusActive {
			allowed[enrollment.UserID] = true
		}
	}
	attendees := make([]string, 0, len(input.AttendeeIDs))
	seen := make(map[string]bool)
	for _, userID := range input.AttendeeIDs {
		if !allowed[userID] {
			return nil, apperrors.NewBadRequestError("attendee must be enrolled in the selected course")
		}
		if !seen[userID] {
			attendees = append(attendees, userID)
			seen[userID] = true
		}
	}

	id, err := app.Repositories.CalendarEvent.Create(ctx, domain.CalendarEvent{
		OwnerID:         requesterID,
		CourseID:        input.CourseID,
		Title:           input.Title,
		Description:     strings.TrimSpace(input.Description),
		StartsAt:        input.StartsAt,
		EndsAt:          input.EndsAt,
		AllDay:          input.AllDay,
		RecurrenceRule:  input.RecurrenceRule,
		ReminderMinutes: input.ReminderMinutes,
		AttendeeIDs:     attendees,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CalendarEventCreateError, err)
	}

	return &CreateOutput{Data: EventData{
		ID:              id,
		OwnerID:         requesterID,
		CourseID:        input.CourseID,
		Title:           input.Title,
		Description:     strings.TrimSpace(input.Description),
		StartsAt:        formatTime(input.StartsAt),
		EndsAt:          formatOptionalTime(input.EndsAt),
		AllDay:          input.AllDay,
		RecurrenceRule:  input.RecurrenceRule,
		ReminderMinutes: input.ReminderMinutes,
		AttendeeIDs:     attendees,
		Source:          "manual",
	}}, nil
}
