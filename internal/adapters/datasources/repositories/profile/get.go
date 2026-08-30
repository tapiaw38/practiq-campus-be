package profile

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.Profile, error) {
	query := `
		SELECT id, profile_type, full_name, avatar_url, bio, created_at, updated_at
		FROM campus_profiles
		WHERE id = $1
	`
	var p domain.Profile
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.ProfileType, &p.FullName, &p.AvatarURL, &p.Bio, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}
