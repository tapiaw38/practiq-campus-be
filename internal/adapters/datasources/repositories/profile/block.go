package profile

import "context"

func (r *repository) SetBlocked(ctx context.Context, id string, blocked bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE campus_profiles SET is_blocked = $2, updated_at = NOW() WHERE id = $1`, id, blocked)
	return err
}
