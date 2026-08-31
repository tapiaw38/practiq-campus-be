package course_group

import "context"

func (r *repository) AddMember(ctx context.Context, groupID, userID string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO course_group_members (group_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, groupID, userID)
	return err
}

func (r *repository) RemoveMember(ctx context.Context, groupID, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM course_group_members WHERE group_id=$1 AND user_id=$2`, groupID, userID)
	return err
}

func (r *repository) IsMember(ctx context.Context, groupID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM course_group_members WHERE group_id=$1 AND user_id=$2)`, groupID, userID).Scan(&exists)
	return exists, err
}
