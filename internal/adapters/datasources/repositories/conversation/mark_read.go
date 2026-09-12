package conversation

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) MarkRead(ctx context.Context, conversationID, userID string) error {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"conversations c", "c.id = conversation_participants.conversation_id", "c.tenant_id", 3)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx,
		"UPDATE conversation_participants SET last_read_at = NOW() WHERE conversation_id = $1 AND user_id = $2 AND "+guard,
		conversationID, userID, tenantID,
	)
	return err
}
