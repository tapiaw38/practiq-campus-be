package message

import (
	"github.com/tapiaw38/practiq-campus-be/internal/adapters/datasources/repositories/conversation"
	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

type MessageData struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`
	SenderID       string `json:"sender_id"`
	Body           string `json:"body"`
	SentAt         string `json:"sent_at"`
}

func toMessageData(m domain.Message) MessageData {
	return MessageData{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
		Body:           m.Body,
		SentAt:         m.SentAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toMessageDataList(messages []domain.Message) []MessageData {
	out := make([]MessageData, 0, len(messages))
	for _, m := range messages {
		out = append(out, toMessageData(m))
	}
	return out
}

type ConversationData struct {
	ID                  string  `json:"id"`
	OtherUserID         string  `json:"other_user_id"`
	OtherUserName       string  `json:"other_user_name"`
	OtherUserEmail      string  `json:"other_user_email"`
	LastMessageBody     string  `json:"last_message_body"`
	LastMessageAt       *string `json:"last_message_at"`
	LastMessageSenderID string  `json:"last_message_sender_id"`
	Unread              bool    `json:"unread"`
}

func toConversationData(s conversation.ConversationSummary, otherUserName, otherUserEmail string) ConversationData {
	data := ConversationData{
		ID:                  s.ID,
		OtherUserID:         s.OtherUserID,
		OtherUserName:       otherUserName,
		OtherUserEmail:      otherUserEmail,
		LastMessageBody:     s.LastMessageBody,
		LastMessageSenderID: s.LastMessageSenderID,
		Unread:              s.Unread,
	}
	if s.LastMessageAt != nil {
		v := s.LastMessageAt.Format("2006-01-02T15:04:05Z")
		data.LastMessageAt = &v
	}
	return data
}
