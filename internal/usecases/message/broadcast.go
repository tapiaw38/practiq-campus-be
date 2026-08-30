package message

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	BroadcastUsecase interface {
		Execute(context.Context, string, bool, string, BroadcastInput) (*BroadcastOutput, apperrors.ApplicationError)
	}

	broadcastUsecase struct {
		contextFactory appcontext.Factory
	}

	BroadcastInput struct {
		Body string
	}

	BroadcastOutput struct {
		Data struct {
			Sent int `json:"sent"`
		} `json:"data"`
	}
)

func NewBroadcastUsecase(contextFactory appcontext.Factory) BroadcastUsecase {
	return &broadcastUsecase{contextFactory: contextFactory}
}

// Execute reuses the plain 1:1 conversation model — one message gets
// appended to (or starts) a separate conversation with every actively
// enrolled student, rather than introducing a group-message concept.
func (u *broadcastUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string, input BroadcastInput) (*BroadcastOutput, apperrors.ApplicationError) {
	if strings.TrimSpace(input.Body) == "" {
		return nil, apperrors.NewBadRequestError("body is required")
	}

	app := u.contextFactory()

	if !isSuperAdmin {
		c, err := app.Repositories.Course.Get(ctx, courseID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
		}
		if c == nil || c.OwnerID != requesterID {
			return nil, apperrors.NewForbiddenError()
		}
	}

	enrollments, err := app.Repositories.Enrollment.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentListError, err)
	}

	sent := 0
	for _, e := range enrollments {
		if e.Status != domain.EnrollmentStatusActive {
			continue
		}
		conversationID, appErr := findOrCreateConversation(ctx, app, requesterID, e.UserID)
		if appErr != nil {
			return nil, appErr
		}
		if _, err := app.Repositories.Conversation.AddMessage(ctx, conversationID, requesterID, input.Body); err != nil {
			return nil, apperrors.NewApplicationError(mappings.MessageCreateError, err)
		}
		sent++
	}

	output := &BroadcastOutput{}
	output.Data.Sent = sent
	return output, nil
}
