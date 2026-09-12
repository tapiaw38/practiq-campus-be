package course_group

import (
	"context"

	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

func (r *repository) AddMember(ctx context.Context, groupID, userID string) error {
	guard, tenantID, err := groupGuard(ctx, "$1", 3)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx,
		`INSERT INTO course_group_members (group_id, user_id)
		 SELECT $1, $2 WHERE `+guard+` ON CONFLICT DO NOTHING`, groupID, userID, tenantID)
	return err
}

func (r *repository) RemoveMember(ctx context.Context, groupID, userID string) error {
	guard, tenantID, err := groupGuard(ctx, "$1", 3)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx,
		`DELETE FROM course_group_members WHERE group_id=$1 AND user_id=$2 AND `+guard,
		groupID, userID, tenantID)
	return err
}

// IsMember gates group-restricted content, so it answers within one
// institution: the same person in another is not a member here.
func (r *repository) IsMember(ctx context.Context, groupID, userID string) (bool, error) {
	tenantID := tenantcontext.ID(ctx)
	if tenantID == "" {
		return false, tenantcontext.ErrNoTenant
	}
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM course_group_members m
			JOIN course_groups g ON g.id = m.group_id
			JOIN courses c ON c.id = g.course_id AND c.tenant_id = $3
			WHERE m.group_id=$1 AND m.user_id=$2
		)`, groupID, userID, tenantID).Scan(&exists)
	return exists, err
}
