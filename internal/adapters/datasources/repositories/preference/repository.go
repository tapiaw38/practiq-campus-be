package preference

import (
	"context"
	"database/sql"
	"encoding/json"
)

type Repository interface {
	Get(context.Context, string, string) (json.RawMessage, bool, error)
	Upsert(context.Context, string, string, json.RawMessage) error
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

func (r *repository) Get(ctx context.Context, userID, scope string) (json.RawMessage, bool, error) {
	var settings json.RawMessage
	err := r.db.QueryRowContext(ctx,
		`SELECT settings FROM user_preferences WHERE user_id = $1 AND scope = $2`, userID, scope,
	).Scan(&settings)
	if err == sql.ErrNoRows {
		return json.RawMessage(`{}`), false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return settings, true, nil
}

func (r *repository) Upsert(ctx context.Context, userID, scope string, settings json.RawMessage) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_preferences (user_id, scope, settings)
		VALUES ($1, $2, $3::jsonb)
		ON CONFLICT (user_id, scope) DO UPDATE
		SET settings = EXCLUDED.settings, updated_at = now()`, userID, scope, settings)
	return err
}
