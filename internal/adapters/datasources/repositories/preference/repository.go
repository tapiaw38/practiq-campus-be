package preference

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

type Repository interface {
	Get(context.Context, string, string) (json.RawMessage, bool, error)
	Upsert(context.Context, string, string, json.RawMessage) error
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

// Preferences belong to a person within one institution. Somebody who teaches
// at two keeps a separate setup in each, so the tenant is part of the key and
// not merely a filter.
func (r *repository) Get(ctx context.Context, userID, scope string) (json.RawMessage, bool, error) {
	tenantID := tenantcontext.ID(ctx)
	if tenantID == "" {
		return nil, false, tenantcontext.ErrNoTenant
	}
	var settings json.RawMessage
	err := r.db.QueryRowContext(ctx,
		`SELECT settings FROM user_preferences WHERE tenant_id = $1 AND user_id = $2 AND scope = $3`,
		tenantID, userID, scope,
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
	tenantID := tenantcontext.ID(ctx)
	if tenantID == "" {
		return tenantcontext.ErrNoTenant
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_preferences (tenant_id, user_id, scope, settings)
		VALUES ($1, $2, $3, $4::jsonb)
		ON CONFLICT (tenant_id, user_id, scope) DO UPDATE
		SET settings = EXCLUDED.settings, updated_at = now()`, tenantID, userID, scope, settings)
	return err
}
