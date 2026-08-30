package domain

import "time"

// Conversation is always a plain 1:1 DM in v1 — a teacher broadcast reuses
// this same shape by fanning a message out into one conversation per
// recipient rather than introducing a separate group-message concept.
type Conversation struct {
	ID        string
	CreatedAt time.Time
}

type Message struct {
	ID             string
	ConversationID string
	SenderID       string
	Body           string
	SentAt         time.Time
}
