package message

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/identity"
)

type (
	ListConversationsUsecase interface {
		Execute(ctx context.Context, requesterID, bearerToken string) (*ListConversationsOutput, apperrors.ApplicationError)
	}

	listConversationsUsecase struct {
		contextFactory appcontext.Factory
	}

	ListConversationsOutput struct {
		Data []ConversationData `json:"data"`
	}
)

func NewListConversationsUsecase(contextFactory appcontext.Factory) ListConversationsUsecase {
	return &listConversationsUsecase{contextFactory: contextFactory}
}

func (u *listConversationsUsecase) Execute(ctx context.Context, requesterID, bearerToken string) (*ListConversationsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	summaries, err := app.Repositories.Conversation.ListMine(ctx, requesterID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ConversationListError, err)
	}

	otherIDs := make([]string, 0, len(summaries))
	for _, s := range summaries {
		otherIDs = append(otherIDs, s.OtherUserID)
	}
	names, err := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, otherIDs)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}

	data := make([]ConversationData, 0, len(summaries))
	for _, s := range summaries {
		info := names[s.OtherUserID]
		data = append(data, toConversationData(s, identity.FullName(info, s.OtherUserID), info.Email))
	}

	return &ListConversationsOutput{Data: data}, nil
}
