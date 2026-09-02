package profile

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Upsert(ctx context.Context, p domain.Profile) error {
	query := `
		INSERT INTO campus_profiles (id, profile_type, avatar_url, bio)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			profile_type = $2,
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.ProfileType, p.AvatarURL, p.Bio)
	return err
}
