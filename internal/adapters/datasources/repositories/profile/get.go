package profile

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

const selectProfileColumns = `id, profile_type, full_name, email, avatar_url, bio, is_blocked, created_at, updated_at`

func scanProfile(row *sql.Row) (*domain.Profile, error) {
	var p domain.Profile
	var email sql.NullString
	err := row.Scan(
		&p.ID, &p.ProfileType, &p.FullName, &email, &p.AvatarURL, &p.Bio, &p.IsBlocked, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	p.Email = email.String
	return &p, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Profile, error) {
	query := `SELECT ` + selectProfileColumns + ` FROM campus_profiles WHERE id = $1`
	return scanProfile(r.db.QueryRowContext(ctx, query, id))
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*domain.Profile, error) {
	query := `SELECT ` + selectProfileColumns + ` FROM campus_profiles WHERE email = $1`
	return scanProfile(r.db.QueryRowContext(ctx, query, email))
}

func scanProfileRows(rows *sql.Rows) ([]domain.Profile, error) {
	profiles := make([]domain.Profile, 0)
	for rows.Next() {
		var p domain.Profile
		var email sql.NullString
		if err := rows.Scan(&p.ID, &p.ProfileType, &p.FullName, &email, &p.AvatarURL, &p.Bio, &p.IsBlocked, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.Email = email.String
		profiles = append(profiles, p)
	}
	return profiles, rows.Err()
}

func (r *repository) ListByType(ctx context.Context, profileType string) ([]domain.Profile, error) {
	query := `SELECT ` + selectProfileColumns + ` FROM campus_profiles WHERE profile_type = $1 ORDER BY full_name`
	rows, err := r.db.QueryContext(ctx, query, profileType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProfileRows(rows)
}

func (r *repository) ListAll(ctx context.Context, search string, limit, offset int) ([]domain.Profile, int, error) {
	filter := "%" + search + "%"
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM campus_profiles WHERE full_name ILIKE $1 OR email ILIKE $1`, filter).Scan(&total); err != nil {
		return nil, 0, err
	}
	query := `SELECT ` + selectProfileColumns + ` FROM campus_profiles WHERE full_name ILIKE $1 OR email ILIKE $1 ORDER BY full_name, email LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, filter, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	profiles, err := scanProfileRows(rows)
	return profiles, total, err
}
