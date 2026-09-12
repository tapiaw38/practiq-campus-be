package course_group

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// ErrNotInTenant is returned when the course or group belongs to another
// institution.
var ErrNotInTenant = errors.New("resource does not belong to this tenant")

// groupGuard constrains a write to a group of the tenant, walking group →
// course.
func groupGuard(ctx context.Context, column string, placeholder int) (string, any, error) {
	return tenantcontext.ExistsIn(ctx,
		"course_groups g JOIN courses c ON c.id = g.course_id", "g.id = "+column, "c.tenant_id", placeholder)
}

func (r *repository) Create(ctx context.Context, g domain.CourseGroup) (string, error) {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx, "courses c", "c.id = $1", "c.tenant_id", 3)
	if err != nil {
		return "", err
	}
	var id string
	err = r.db.QueryRowContext(ctx,
		`INSERT INTO course_groups (course_id, name) SELECT $1, $2 WHERE `+guard+` RETURNING id`,
		g.CourseID, g.Name, tenantID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotInTenant
	}
	return id, err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	guard, tenantID, err := tenantcontext.ExistsIn(ctx,
		"courses c", "c.id = course_groups.course_id", "c.tenant_id", 2)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `DELETE FROM course_groups WHERE id=$1 AND `+guard, id, tenantID)
	return err
}
