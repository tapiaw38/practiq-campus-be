package conversation

import (
	"context"
	"database/sql"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type ConversationSummary struct {
	ID                  string
	OtherUserID         string
	LastMessageBody     string
	LastMessageAt       *time.Time
	LastMessageSenderID string
	Unread              bool
}

type Repository interface {
	// FindDirectBetween returns the existing 1:1 conversation between
	// exactly these two users, or "" if none exists yet.
	FindDirectBetween(ctx context.Context, userA, userB string) (string, error)
	CreateDirect(ctx context.Context, userA, userB string) (string, error)
	AddMessage(ctx context.Context, conversationID, senderID, body string) (string, error)
	GetMessage(ctx context.Context, id string) (*domain.Message, error)
	ListMessages(ctx context.Context, conversationID string) ([]domain.Message, error)
	IsParticipant(ctx context.Context, conversationID, userID string) (bool, error)
	ListMine(ctx context.Context, userID string) ([]ConversationSummary, error)
	MarkRead(ctx context.Context, conversationID, userID string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
