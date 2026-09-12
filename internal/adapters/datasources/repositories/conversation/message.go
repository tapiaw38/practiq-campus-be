package conversation

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// ErrConversationNotInTenant is returned when a conversation exists but is
// held in another institution.
var ErrConversationNotInTenant = errors.New("conversation does not belong to this tenant")

const selectMessageColumns = "m.id, m.conversation_id, m.sender_id, m.body, m.sent_at"

func (r *repository) GetMessage(ctx context.Context, id string) (*domain.Message, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("conversations", "c", "m.conversation_id").
		Where("m.id = ?", id).
		SQL(selectMessageColumns, "messages m")
	if err != nil {
		return nil, err
	}

	var m domain.Message
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Body, &m.SentAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// AddMessage refuses a conversation from another institution. The insert
// selects its conversation through the tenant, so a conversation id guessed
// or carried over from another context writes nothing.
func (r *repository) AddMessage(ctx context.Context, conversationID, senderID, body string) (string, error) {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"conversations c", "c.id = $1", "c.tenant_id", 4)
	if err != nil {
		return "", err
	}

	var id string
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO messages (conversation_id, sender_id, body)
		SELECT $1, $2, $3 WHERE `+guard+`
		RETURNING id`, conversationID, senderID, body, tenantID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrConversationNotInTenant
	}
	if err != nil {
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
	query, args, err := tenantcontext.NewQuery(ctx).
		Through("conversations", "c", "m.conversation_id").
		Where("m.conversation_id = ?", conversationID).
		SQL(selectMessageColumns, "messages m")
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query+" ORDER BY m.sent_at ASC", args...)
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

// IsParticipant is the guard the whole of messaging authorises through, so it
// answers for one institution at a time.
//
// Membership of a conversation does not travel: the same person, in another
// institution, is not a participant of a thread held in this one. Without the
// scope, selecting a different tenant kept the old conversation readable.
func (r *repository) IsParticipant(ctx context.Context, conversationID, userID string) (bool, error) {
	tenantID := tenantcontext.ID(ctx)
	if tenantID == "" {
		return false, tenantcontext.ErrNoTenant
	}
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM conversation_participants cp
			JOIN conversations c ON c.id = cp.conversation_id AND c.tenant_id = $3
			WHERE cp.conversation_id = $1 AND cp.user_id = $2
		)`, conversationID, userID, tenantID).Scan(&exists)
	return exists, err
}
