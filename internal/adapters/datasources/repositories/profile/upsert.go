package profile

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

func (r *repository) Upsert(ctx context.Context, p domain.Profile) error {
	query := `
		INSERT INTO campus_profiles (id, profile_type, full_name, email, avatar_url, bio)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			profile_type = $2,
			full_name = $3,
			email = COALESCE(NULLIF($4, ''), campus_profiles.email),
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.ProfileType, p.FullName, p.Email, p.AvatarURL, p.Bio)
	return err
}
