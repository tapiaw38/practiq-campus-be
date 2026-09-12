package tenant

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
)

const columns = "id, school_id, status, activated_at, deactivated_at, created_at"

func scan(row interface{ Scan(...any) error }) (*domain.Tenant, error) {
	var value domain.Tenant
	if err := row.Scan(&value.ID, &value.SchoolID, &value.Status, &value.ActivatedAt, &value.DeactivatedAt, &value.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &value, nil
}

func (r *repository) Create(ctx context.Context, value domain.Tenant) (*domain.Tenant, error) {
	return scan(r.db.QueryRowContext(ctx, `
		INSERT INTO campus_tenants (school_id, status) VALUES ($1, $2)
		ON CONFLICT (school_id) DO UPDATE SET status = EXCLUDED.status,
		deactivated_at = CASE WHEN EXCLUDED.status = 'active' THEN NULL ELSE now() END
		RETURNING `+columns, value.SchoolID, value.Status))
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Tenant, error) {
	return scan(r.db.QueryRowContext(ctx, "SELECT "+columns+" FROM campus_tenants WHERE id = $1", id))
}

func (r *repository) GetBySchoolID(ctx context.Context, schoolID string) (*domain.Tenant, error) {
	return scan(r.db.QueryRowContext(ctx, "SELECT "+columns+" FROM campus_tenants WHERE school_id = $1", schoolID))
}

func (r *repository) List(ctx context.Context) ([]domain.Tenant, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT "+columns+" FROM campus_tenants ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []domain.Tenant{}
	for rows.Next() {
		value, err := scan(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *value)
	}
	return result, rows.Err()
}

func (r *repository) SetStatus(ctx context.Context, id, status string) (*domain.Tenant, error) {
	return scan(r.db.QueryRowContext(ctx, `
		UPDATE campus_tenants
		SET status = $2, deactivated_at = CASE WHEN $2 = 'active' THEN NULL ELSE now() END
		WHERE id = $1 RETURNING `+columns, id, status))
}
