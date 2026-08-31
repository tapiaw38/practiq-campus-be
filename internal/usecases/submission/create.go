package submission

import (
	"context"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/unlock"
)

type (
	CreateUsecase interface {
		Execute(context.Context, string, string, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		Content string
	}

	CreateOutput struct {
		Data SubmissionData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

// Execute is any authenticated user — the enrollment check below (not a
// teacherOnly gate) is what actually decides whether they may submit.
func (u *createUsecase) Execute(ctx context.Context, requesterID, assignmentID string, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	if strings.TrimSpace(input.Content) == "" {
		return nil, apperrors.NewBadRequestError("content is required")
	}

	app := u.contextFactory()

	a, err := app.Repositories.Assignment.Get(ctx, assignmentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentGetError, err)
	}
	if a == nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentNotFoundError, nil)
	}

	enrollment, err := app.Repositories.Enrollment.GetByCourseAndUser(ctx, a.CourseID, requesterID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentGetError, err)
	}
	if enrollment == nil {
		return nil, apperrors.NewApplicationError(mappings.SubmissionNotEnrolledError, nil)
	}
	if a.VisibleGroupID != nil {
		member, err := app.Repositories.CourseGroup.IsMember(ctx, *a.VisibleGroupID, requesterID)
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		if !member {
			return nil, apperrors.NewForbiddenError()
		}
	}
	if status, err := unlock.Check(ctx, app.Repositories, a.UnlockAfterType, a.UnlockAfterID, requesterID); err != nil {
		return nil, apperrors.NewInternalError(err)
	} else if status.Locked {
		return nil, apperrors.NewBadRequestError(status.Reason)
	}

	if existing, err := app.Repositories.Submission.GetByAssignmentAndUser(ctx, assignmentID, requesterID); err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubmissionGetError, err)
	} else if existing != nil {
		if existing.GradedAt != nil {
			return nil, apperrors.NewBadRequestError("la entrega ya fue corregida y no puede modificarse")
		}
		if a.DueAt != nil && time.Now().After(*a.DueAt) {
			return nil, apperrors.NewBadRequestError("el plazo de entrega ya venció")
		}
		if err := app.Repositories.Submission.Resubmit(ctx, existing.ID, input.Content); err != nil {
			return nil, apperrors.NewApplicationError(mappings.SubmissionCreateError, err)
		}
		updated, err := app.Repositories.Submission.Get(ctx, existing.ID)
		if err != nil || updated == nil {
			return nil, apperrors.NewApplicationError(mappings.SubmissionGetError, err)
		}
		return &CreateOutput{Data: withAttachments(app, toSubmissionData(*updated))}, nil
	}

	id, err := app.Repositories.Submission.Create(ctx, domain.Submission{
		AssignmentID: assignmentID,
		UserID:       requesterID,
		Content:      input.Content,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubmissionCreateError, err)
	}

	created, err := app.Repositories.Submission.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubmissionGetError, err)
	}
	if created == nil {
		return nil, apperrors.NewInternalError(nil)
	}
	if course, err := app.Repositories.Course.Get(ctx, a.CourseID); err == nil && course != nil {
		_ = app.Repositories.Notification.Create(ctx, domain.Notification{UserID: course.OwnerID, Type: "submission_created", Title: "Nueva entrega", Body: a.Title, Data: `{"submission_id":"` + created.ID + `","assignment_id":"` + assignmentID + `"}`})
	}

	return &CreateOutput{Data: withAttachments(app, toSubmissionData(*created))}, nil
}
