package conversation

import "context"

func (r *repository) MarkRead(ctx context.Context, conversationID, userID string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE conversation_participants SET last_read_at = NOW() WHERE conversation_id = $1 AND user_id = $2",
		conversationID, userID,
	)
	return err
}
