package domain

import "time"

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
