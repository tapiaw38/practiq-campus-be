package message

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-campus-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/errors/mappings"
)

type (
	ListConversationsUsecase interface {
		Execute(context.Context, string) (*ListConversationsOutput, apperrors.ApplicationError)
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

func (u *listConversationsUsecase) Execute(ctx context.Context, requesterID string) (*ListConversationsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	summaries, err := app.Repositories.Conversation.ListMine(ctx, requesterID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ConversationListError, err)
	}

	data := make([]ConversationData, 0, len(summaries))
	for _, s := range summaries {
		otherProfile, err := app.Repositories.Profile.Get(ctx, s.OtherUserID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
		}
		data = append(data, toConversationData(s, otherProfile))
	}

	return &ListConversationsOutput{Data: data}, nil
}
