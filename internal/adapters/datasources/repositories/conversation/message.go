package conversation

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) GetMessage(ctx context.Context, id string) (*domain.Message, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, conversation_id, sender_id, body, sent_at FROM messages WHERE id = $1", id)

	var m domain.Message
	err := row.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Body, &m.SentAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) AddMessage(ctx context.Context, conversationID, senderID, body string) (string, error) {
	query := `
		INSERT INTO messages (conversation_id, sender_id, body)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id string
	if err := r.db.QueryRowContext(ctx, query, conversationID, senderID, body).Scan(&id); err != nil {
		return "", err
	}

	// Sending a message counts as having read up to that point — otherwise
	// the sender's own conversation would immediately show as unread to
	// themselves.
	if _, err := r.db.ExecContext(ctx,
		"UPDATE conversation_participants SET last_read_at = NOW() WHERE conversation_id = $1 AND user_id = $2",
		conversationID, senderID,
	); err != nil {
		return "", err
	}

	return id, nil
}

func (r *repository) ListMessages(ctx context.Context, conversationID string) ([]domain.Message, error) {
	query := `
		SELECT id, conversation_id, sender_id, body, sent_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY sent_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Body, &m.SentAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func (r *repository) IsParticipant(ctx context.Context, conversationID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM conversation_participants WHERE conversation_id = $1 AND user_id = $2)",
		conversationID, userID,
	).Scan(&exists)
	return exists, err
}
