package conversation

import (
	"context"
	"database/sql"
)

func (r *repository) ListMine(ctx context.Context, userID string) ([]ConversationSummary, error) {
	query := `
		SELECT
			c.id,
			(
				SELECT cp2.user_id FROM conversation_participants cp2
				WHERE cp2.conversation_id = c.id AND cp2.user_id != $1
				LIMIT 1
			) AS other_user_id,
			m.body,
			m.sent_at,
			m.sender_id,
			(
				m.sent_at IS NOT NULL
				AND m.sender_id != $1
				AND (cp.last_read_at IS NULL OR m.sent_at > cp.last_read_at)
			) AS is_unread
		FROM conversations c
		JOIN conversation_participants cp ON cp.conversation_id = c.id AND cp.user_id = $1
		LEFT JOIN LATERAL (
			SELECT body, sent_at, sender_id FROM messages WHERE conversation_id = c.id ORDER BY sent_at DESC LIMIT 1
		) m ON true
		ORDER BY m.sent_at DESC NULLS LAST
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []ConversationSummary
	for rows.Next() {
		var s ConversationSummary
		var lastBody, lastSender sql.NullString
		var lastAt sql.NullTime
		if err := rows.Scan(&s.ID, &s.OtherUserID, &lastBody, &lastAt, &lastSender, &s.Unread); err != nil {
			return nil, err
		}
		s.LastMessageBody = lastBody.String
		s.LastMessageSenderID = lastSender.String
		if lastAt.Valid {
			s.LastMessageAt = &lastAt.Time
		}
		summaries = append(summaries, s)
	}
	return summaries, rows.Err()
}
