package message

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	ReplyUsecase interface {
		Execute(context.Context, string, string, ReplyInput) (*SendOutput, apperrors.ApplicationError)
	}

	replyUsecase struct {
		contextFactory appcontext.Factory
	}

	ReplyInput struct {
		Body string
	}
)

func NewReplyUsecase(contextFactory appcontext.Factory) ReplyUsecase {
	return &replyUsecase{contextFactory: contextFactory}
}

func (u *replyUsecase) Execute(ctx context.Context, requesterID, conversationID string, input ReplyInput) (*SendOutput, apperrors.ApplicationError) {
	if strings.TrimSpace(input.Body) == "" {
		return nil, apperrors.NewBadRequestError("body is required")
	}

	app := u.contextFactory()

	isParticipant, err := app.Repositories.Conversation.IsParticipant(ctx, conversationID, requesterID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ConversationListError, err)
	}
	if !isParticipant {
		return nil, apperrors.NewForbiddenError()
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
	return &SendOutput{Data: toMessageData(*created)}, nil
}
