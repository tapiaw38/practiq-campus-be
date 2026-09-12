package conversation

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// FindDirectBetween looks only inside the tenant in context.
//
// Two people can share more than one institution, and a conversation belongs
// to exactly one of them. Without the scope, writing to a colleague from one
// institution would reopen the thread held in the other and show its history
// to a context it never belonged to.
func (r *repository) FindDirectBetween(ctx context.Context, userA, userB string) (string, error) {
	query, args, err := tenantcontext.NewQuery(ctx).
		Own("c.tenant_id").
		Where("cp1.user_id = ?", userA).
		Where("cp2.user_id = ?", userB).
		Where(`(SELECT COUNT(*) FROM conversation_participants cp3
		        WHERE cp3.conversation_id = c.id) = 2`).
		SQL("c.id", `conversations c
			JOIN conversation_participants cp1 ON cp1.conversation_id = c.id
			JOIN conversation_participants cp2 ON cp2.conversation_id = c.id`)
	if err != nil {
		return "", err
	}

	var id string
	err = r.db.QueryRowContext(ctx, query+" LIMIT 1", args...).Scan(&id)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return id, err
}

// CreateDirect stamps the conversation with the tenant it is being held in.
func (r *repository) CreateDirect(ctx context.Context, userA, userB string) (string, error) {
	tenantID := tenantcontext.ID(ctx)
	if tenantID == "" {
		return "", tenantcontext.ErrNoTenant
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var id string
	if err := tx.QueryRowContext(ctx,
		"INSERT INTO conversations (tenant_id) VALUES ($1) RETURNING id", tenantID).Scan(&id); err != nil {
		return "", err
	}

	if _, err := tx.ExecContext(ctx,
		"INSERT INTO conversation_participants (conversation_id, user_id) VALUES ($1, $2), ($1, $3)",
		id, userA, userB,
	); err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	return id, nil
}
