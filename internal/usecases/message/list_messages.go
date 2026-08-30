package message

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	ListMessagesUsecase interface {
		Execute(context.Context, string, string) (*ListMessagesOutput, apperrors.ApplicationError)
	}

	listMessagesUsecase struct {
		contextFactory appcontext.Factory
	}

	ListMessagesOutput struct {
		Data []MessageData `json:"data"`
	}
)

func NewListMessagesUsecase(contextFactory appcontext.Factory) ListMessagesUsecase {
	return &listMessagesUsecase{contextFactory: contextFactory}
}

func (u *listMessagesUsecase) Execute(ctx context.Context, requesterID, conversationID string) (*ListMessagesOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	isParticipant, err := app.Repositories.Conversation.IsParticipant(ctx, conversationID, requesterID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ConversationListError, err)
	}
	if !isParticipant {
		return nil, apperrors.NewForbiddenError()
	}

	messages, err := app.Repositories.Conversation.ListMessages(ctx, conversationID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MessageListError, err)
	}

	return &ListMessagesOutput{Data: toMessageDataList(messages)}, nil
}
