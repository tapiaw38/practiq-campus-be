package message

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	MarkReadUsecase interface {
		Execute(context.Context, string, string) apperrors.ApplicationError
	}

	markReadUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewMarkReadUsecase(contextFactory appcontext.Factory) MarkReadUsecase {
	return &markReadUsecase{contextFactory: contextFactory}
}

func (u *markReadUsecase) Execute(ctx context.Context, requesterID, conversationID string) apperrors.ApplicationError {
	app := u.contextFactory()

	isParticipant, err := app.Repositories.Conversation.IsParticipant(ctx, conversationID, requesterID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.ConversationListError, err)
	}
	if !isParticipant {
		return apperrors.NewForbiddenError()
	}

	if err := app.Repositories.Conversation.MarkRead(ctx, conversationID, requesterID); err != nil {
		return apperrors.NewApplicationError(mappings.ConversationCreateError, err)
	}
	return nil
}
