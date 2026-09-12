package submission

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// ErrAssignmentNotInTenant is returned when the assignment belongs to another
// institution's course.
var ErrAssignmentNotInTenant = errors.New("assignment does not belong to this tenant")

// assignmentGuard constrains a write to an assignment of the tenant in context.
func assignmentGuard(ctx context.Context, column string, placeholder int) (string, any, error) {
	return tenantcontext.ExistsIn(ctx,
		"assignments a JOIN courses c ON c.id = a.course_id", "a.id = "+column, "c.tenant_id", placeholder)
}

func (r *repository) Create(ctx context.Context, s domain.Submission) (string, error) {
	guard, tenantID, err := assignmentGuard(ctx, "$1", 5)
	if err != nil {
		return "", err
	}
	var id string
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO submissions (assignment_id, user_id, content, status)
		SELECT $1, $2, $3, $4 WHERE `+guard+`
		RETURNING id`,
		s.AssignmentID, s.UserID, s.Content, domain.SubmissionStatusSubmitted, tenantID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrAssignmentNotInTenant
	}
	return id, err
}
