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
	SendUsecase interface {
		Execute(context.Context, string, SendInput) (*SendOutput, apperrors.ApplicationError)
	}

	sendUsecase struct {
		contextFactory appcontext.Factory
	}

	SendInput struct {
		// ToUserID takes priority when set — the compose UI always resolves
		// a concrete person via search first, so it should never need to
		// re-derive identity from an email string. ToEmail stays as a
		// fallback for any other caller.
		ToUserID    string
		ToEmail     string
		Body        string
		BearerToken string
	}

	SendOutput struct {
		Data MessageData `json:"data"`
	}
)

func NewSendUsecase(contextFactory appcontext.Factory) SendUsecase {
	return &sendUsecase{contextFactory: contextFactory}
}

// findOrCreateConversation is shared by Send and Broadcast — both ultimately
// just need "the 1:1 conversation between these two people," created on
// first contact.
func findOrCreateConversation(ctx context.Context, app *appcontext.Context, userA, userB string) (string, apperrors.ApplicationError) {
	id, err := app.Repositories.Conversation.FindDirectBetween(ctx, userA, userB)
	if err != nil {
		return "", apperrors.NewApplicationError(mappings.ConversationListError, err)
	}
	if id != "" {
		return id, nil
	}

	id, err = app.Repositories.Conversation.CreateDirect(ctx, userA, userB)
	if err != nil {
		return "", apperrors.NewApplicationError(mappings.ConversationCreateError, err)
	}
	return id, nil
}

func (u *sendUsecase) Execute(ctx context.Context, requesterID string, input SendInput) (*SendOutput, apperrors.ApplicationError) {
	if strings.TrimSpace(input.Body) == "" {
		return nil, apperrors.NewBadRequestError("body is required")
	}

	app := u.contextFactory()

	var recipientID string
	if input.ToUserID != "" {
		recipient, err := app.Repositories.Profile.Get(ctx, input.ToUserID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
		}
		if recipient == nil {
			return nil, apperrors.NewApplicationError(mappings.MessageRecipientNotFoundError, nil)
		}
		recipientID = recipient.ID
	} else {
		authUser, err := app.Integrations.AuthAPI.GetByEmail(ctx, input.BearerToken, input.ToEmail)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
		}
		if authUser == nil {
			return nil, apperrors.NewApplicationError(mappings.MessageRecipientNotFoundError, nil)
		}
		recipient, err := app.Repositories.Profile.Get(ctx, authUser.Username)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
		}
		if recipient == nil {
			return nil, apperrors.NewApplicationError(mappings.MessageRecipientNotFoundError, nil)
		}
		recipientID = recipient.ID
	}

	shares, err := app.Repositories.Enrollment.SharesCourseWith(ctx, requesterID, recipientID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentGetError, err)
	}
	if !shares {
		return nil, apperrors.NewApplicationError(mappings.MessageNoSharedCourseError, nil)
	}

	conversationID, appErr := findOrCreateConversation(ctx, app, requesterID, recipientID)
	if appErr != nil {
		return nil, appErr
	}

	msgID, err := app.Repositories.Conversation.AddMessage(ctx, conversationID, requesterID, input.Body)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MessageCreateError, err)
	}

	created, err := app.Repositories.Conversation.GetMessage(ctx, msgID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MessageListError, err)
	}
	if created == nil {
		return nil, apperrors.NewInternalError(nil)
	}
	_ = app.Repositories.Notification.Create(ctx, domain.Notification{UserID: recipientID, Type: "message", Title: "Nuevo mensaje", Body: input.Body, Data: `{"conversation_id":"` + conversationID + `"}`})
	return &SendOutput{Data: toMessageData(*created)}, nil
}
