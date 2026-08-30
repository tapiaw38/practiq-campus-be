package conversation

import (
	"context"
	"database/sql"
)

func (r *repository) FindDirectBetween(ctx context.Context, userA, userB string) (string, error) {
	query := `
		SELECT cp1.conversation_id
		FROM conversation_participants cp1
		JOIN conversation_participants cp2
			ON cp2.conversation_id = cp1.conversation_id AND cp2.user_id = $2
		WHERE cp1.user_id = $1
		AND (
			SELECT COUNT(*) FROM conversation_participants cp3
			WHERE cp3.conversation_id = cp1.conversation_id
		) = 2
		LIMIT 1
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, userA, userB).Scan(&id)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return id, err
}

func (r *repository) CreateDirect(ctx context.Context, userA, userB string) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var id string
	if err := tx.QueryRowContext(ctx, "INSERT INTO conversations DEFAULT VALUES RETURNING id").Scan(&id); err != nil {
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
