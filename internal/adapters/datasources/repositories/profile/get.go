package profile

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

const selectProfileColumns = `id, profile_type, avatar_url, bio, is_blocked, created_at, updated_at`

func scanProfile(row *sql.Row) (*domain.Profile, error) {
	var p domain.Profile
	err := row.Scan(&p.ID, &p.ProfileType, &p.AvatarURL, &p.Bio, &p.IsBlocked, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Profile, error) {
	query := `SELECT ` + selectProfileColumns + ` FROM campus_profiles WHERE id = $1`
	return scanProfile(r.db.QueryRowContext(ctx, query, id))
}

func scanProfileRows(rows *sql.Rows) ([]domain.Profile, error) {
	profiles := make([]domain.Profile, 0)
	for rows.Next() {
		var p domain.Profile
		if err := rows.Scan(&p.ID, &p.ProfileType, &p.AvatarURL, &p.Bio, &p.IsBlocked, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, rows.Err()
}

func (r *repository) ListByType(ctx context.Context, profileType string) ([]domain.Profile, error) {
	query := `SELECT ` + selectProfileColumns + ` FROM campus_profiles WHERE profile_type = $1 ORDER BY created_at`
	rows, err := r.db.QueryContext(ctx, query, profileType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProfileRows(rows)
}

func (r *repository) ListAll(ctx context.Context, limit, offset int) ([]domain.Profile, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM campus_profiles`).Scan(&total); err != nil {
		return nil, 0, err
	}
	query := `SELECT ` + selectProfileColumns + ` FROM campus_profiles ORDER BY created_at LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	profiles, err := scanProfileRows(rows)
	return profiles, total, err
}
